// Package repository
package repository

import (
	"context"
	"database/sql"
	"time"

	"user-auth-app/internal/domain"
	"user-auth-app/internal/repository/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var dbQueryDuration = promauto.NewHistogram(
	prometheus.HistogramOpts{Name: "db_query_duration_seconds", Help: "DB query latency"},
)

type UserRepository struct {
	db *sqlc.Queries
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: sqlc.New(pool)}
}

func (r *UserRepository) CreateUser(ctx context.Context, user domain.User, passwordHash string) (domain.User, error) {
	start := time.Now()
	defer func() { dbQueryDuration.Observe(time.Since(start).Seconds()) }()
	created, err := r.db.CreateUser(ctx, sqlc.CreateUserParams{
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: passwordHash,
		Role:         user.Role,
	})
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		ID:        created.ID,
		Username:  created.Username,
		Email:     created.Email,
		Role:      created.Role,
		CreatedAt: created.CreatedAt,
	}, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, string, error) {
	start := time.Now()
	defer func() { dbQueryDuration.Observe(time.Since(start).Seconds()) }()
	u, err := r.db.GetUserByEmail(ctx, email)
	if err == sql.ErrNoRows {
		return domain.User{}, "", nil // Not found
	} else if err != nil {
		return domain.User{}, "", err
	}
	return domain.User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
	}, u.PasswordHash, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int32) (domain.User, error) {
	start := time.Now()
	defer func() { dbQueryDuration.Observe(time.Since(start).Seconds()) }()
	u, err := r.db.GetUserByID(ctx, id)
	if err == sql.ErrNoRows {
		return domain.User{}, nil
	} else if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
	}, nil
}
