// Package services
package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"user-auth-app/internal/domain"
	"user-auth-app/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo      *repository.UserRepository
	logger    *zerolog.Logger
	jwtSecret string
	redis     *redis.Client
	cache     sync.Map // Fallback in-memory cache
	cacheTTL  time.Duration
	natsConn  *nats.Conn
}

func NewAuthService(repo *repository.UserRepository, logger *zerolog.Logger, jwtSecret string, redisURL string, natsURL string) *AuthService {
	var rdb *redis.Client
	if redisURL != "" {
		rdb = redis.NewClient(&redis.Options{Addr: redisURL})
		if err := rdb.Ping(context.Background()).Err(); err != nil {
			logger.Warn().Err(err).Msg("Redis unavailable, using sync.Map fallback")
			rdb = nil
		}
	}
	nc, err := nats.Connect(natsURL)
	if err != nil {
		logger.Warn().Err(err).Msg("NATS unavailable, skipping async tasks")
	}
	return &AuthService{
		repo:      repo,
		logger:    logger,
		jwtSecret: jwtSecret,
		redis:     rdb,
		cacheTTL:  5 * time.Minute,
		natsConn:  nc,
	}
}

func (s *AuthService) Register(ctx context.Context, username, email, password, role string) (domain.User, error) {
	if password == "" {
		return domain.User{}, errors.New("password required")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, err
	}
	user := domain.User{Username: username, Email: email, Role: role}
	created, err := s.repo.CreateUser(ctx, user, string(hash))
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to create user")
		return domain.User{}, err
	}

	// Publish to NATS for async email verification
	if s.natsConn != nil {
		msg := fmt.Sprintf(`{"user_id": %d, "email": "%s"}`, created.ID, created.Email)
		if err := s.natsConn.Publish("user.verify", []byte(msg)); err != nil {
			s.logger.Error().Err(err).Msg("failed to publish verification message")
		}
	}

	return created, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, hash, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get user")
		return "", err
	}
	if user.ID == 0 {
		return "", errors.New("user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) GetProfile(ctx context.Context, userID int32) (domain.User, error) {
	cacheKey := fmt.Sprintf("user:%d", userID)

	// Try cache first
	if s.redis != nil {
		val, err := s.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var user domain.User
			if err := json.Unmarshal([]byte(val), &user); err == nil {
				return user, nil
			}
		} else if err != redis.Nil {
			s.logger.Error().Err(err).Msg("Redis get failed")
		}
	} else {
		if v, ok := s.cache.Load(cacheKey); ok {
			return v.(domain.User), nil
		}
	}

	// Fetch from DB
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get profile")
		return domain.User{}, err
	}

	// Cache it
	if s.redis != nil {
		userJSON, _ := json.Marshal(user)
		if err := s.redis.Set(ctx, cacheKey, userJSON, s.cacheTTL).Err(); err != nil {
			s.logger.Error().Err(err).Msg("Redis set failed")
		}
	} else {
		s.cache.Store(cacheKey, user)
		// Simulate TTL with goroutine (not production-ready; use for dev)
		go func() {
			time.Sleep(s.cacheTTL)
			s.cache.Delete(cacheKey)
		}()
	}

	return user, nil
}
