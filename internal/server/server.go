// Package server
package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"user-auth-app/internal/config"
	"user-auth-app/internal/handlers"
	"user-auth-app/internal/middleware"
	"user-auth-app/internal/repository"
	"user-auth-app/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	requestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{Name: "http_requests_total", Help: "Total requests"},
		[]string{"path", "method"},
	)
	requestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{Name: "http_request_duration_seconds", Help: "Request latency"},
		[]string{"path", "method"},
	)
)

type Server struct {
	httpServer *http.Server
	pool       *pgxpool.Pool
}

func NewServer(cfg *config.Config) *Server {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger := zerolog.New(output).With().Timestamp().Logger()

	logLevel, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		logLevel = zerolog.InfoLevel
		logger.Warn().Str("log_level", cfg.LogLevel).Msg("Invalid log level, defaulting to Info")
	}
	zerolog.SetGlobalLevel(logLevel)
	pool, err := pgxpool.New(context.Background(), cfg.DBURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to DB")
	}

	userRepo := repository.NewUserRepository(pool)
	authService := services.NewAuthService(userRepo, &logger, cfg.JWTSecret, cfg.RedisURL, cfg.NatsURL)
	authHandler := handlers.NewAuthHandler(authService, &logger, cfg.Timeout)

	r := chi.NewRouter()

	// Metrics middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			path := r.URL.Path
			method := r.Method
			defer func() {
				duration := time.Since(start).Seconds()
				requestDuration.WithLabelValues(path, method).Observe(duration)
				requestsTotal.WithLabelValues(path, method).Inc()
			}()
			next.ServeHTTP(w, r)
		})
	})

	r.Post("/register", authHandler.Register)
	r.Post("/login", authHandler.Login)
	r.Route("/users", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(cfg.JWTSecret, &logger))
		r.Get("/{id}", authHandler.GetProfile)
	})
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	r.Get("/metrics", promhttp.Handler().ServeHTTP)

	httpServer := &http.Server{
		Addr:    cfg.Port,
		Handler: r,
	}

	return &Server{httpServer: httpServer, pool: pool}
}

func (s *Server) Start() {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server failed")
		}
	}()
	log.Info().Msgf("Server running on %s", s.httpServer.Addr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	}
	s.pool.Close()
	log.Info().Msg("Server exited")
}
