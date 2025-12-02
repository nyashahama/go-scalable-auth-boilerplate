// Package config
package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	DBURL          string
	JWTSecret      string
	Port           string
	LogLevel       string
	Timeout        time.Duration
	RedisURL       string
	NatsURL        string
	Environment    string
	AllowedOrigins []string
	RateLimitRPS   int
	RateLimitBurst int
	JWTExpiry      time.Duration
}

func Load() *Config {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL env var required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET env var required")
	}

	port := getEnvOrDefault("PORT", "8080")
	logLevel := getEnvOrDefault("LOG_LEVEL", "info")
	environment := getEnvOrDefault("ENVIRONMENT", "development")

	timeout := getEnvAsInt("TIMEOUT_SECONDS", 30)
	rateLimitRPS := getEnvAsInt("RATE_LIMIT_RPS", 10)
	rateLimitBurst := getEnvAsInt("RATE_LIMIT_BURST", 20)
	jwtExpiryHours := getEnvAsInt("JWT_EXPIRY_HOURS", 24)

	redisURL := getEnvOrDefault("REDIS_URL", "redis://localhost:6379")
	natsURL := getEnvOrDefault("NATS_URL", "nats://localhost:4222")

	// Parse allowed origins
	originsStr := getEnvOrDefault("ALLOWED_ORIGINS", "*")
	allowedOrigins := strings.Split(originsStr, ",")
	for i, origin := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(origin)
	}

	return &Config{
		DBURL:          dbURL,
		JWTSecret:      jwtSecret,
		Port:           ":" + port,
		LogLevel:       logLevel,
		Timeout:        time.Duration(timeout) * time.Second,
		RedisURL:       redisURL,
		NatsURL:        natsURL,
		Environment:    environment,
		AllowedOrigins: allowedOrigins,
		RateLimitRPS:   rateLimitRPS,
		RateLimitBurst: rateLimitBurst,
		JWTExpiry:      time.Duration(jwtExpiryHours) * time.Hour,
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
