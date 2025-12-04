# User Authentication API

A production-ready authentication API built with Go, featuring JWT authentication, Redis caching, NATS messaging, and comprehensive observability.

## Features

- 🔐 **JWT Authentication** - Secure token-based authentication
- 👤 **User Management** - Registration, login, profile management
- 🚀 **High Performance** - Redis caching for optimized data access
- 📨 **Event-Driven** - NATS messaging for async operations
- 📊 **Observability** - Prometheus metrics, structured logging
- 🛡️ **Security** - Rate limiting, CORS, input validation
- 🏗️ **Clean Architecture** - Separation of concerns, dependency injection
- 🐳 **Docker Ready** - Full Docker Compose setup
- ✅ **Production Ready** - Health checks, graceful shutdown, error handling

## Architecture

```
cmd/api/              # Application entry point
internal/
├── app/              # Application initialization & DI
├── config/           # Configuration management
├── domain/           # Domain models and errors
├── repository/       # Data access layer
├── service/          # Business logic layer
├── handler/          # HTTP handlers
├── middleware/       # HTTP middleware
├── cache/            # Caching abstraction
├── messaging/        # Message broker abstraction
└── validator/        # Input validation
```

## Prerequisites

- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 16
- Redis 7
- NATS 2
- sqlc (for code generation)
- golang-migrate (for migrations)

## Quick Start

### 1. Clone and Setup

```bash
git clone https://github.com/nyashahama/go-scalable-auth-boilerplate
cd user-auth-app
cp .env.example .env
```

### 2. Generate JWT Secret

```bash
make generate-jwt
# Copy output to .env JWT_SECRET
```

### 3. Start Services

```bash
# Start all services with Docker Compose
make docker-up

# View logs
make docker-logs
```

### 4. Run Migrations

```bash
# Create migration
make migrate-create name=create_users_table

# Run migrations
make migrate-up DB_URL="postgres://postgres:admin@localhost:5432/dbname?sslmode=disable"
```

### 5. Generate Database Code

```bash
make sqlc
```

### 6. Build and Run

```bash
# Build binary
make build

# Run application
make run

# Or run with Docker
make docker-build
```

## API Endpoints

### Public Endpoints

#### Register User

```bash
POST /api/v1/register
Content-Type: application/json

{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "securepassword123",
  "role": "user"
}
```

#### Login

```bash
POST /api/v1/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "securepassword123"
}
```

### Protected Endpoints (Require Bearer Token)

#### Get User Profile

```bash
GET /api/v1/users/{id}
Authorization: Bearer <token>
```

#### Refresh Token

```bash
POST /api/v1/auth/refresh
Authorization: Bearer <token>
```

### Health & Monitoring

```bash
GET /health      # Comprehensive health check
GET /ready       # Readiness probe
GET /live        # Liveness probe
GET /metrics     # Prometheus metrics
```

## Configuration

All configuration is done via environment variables. See `.env.example` for all available options.

### Key Configuration Options

| Variable           | Description                                    | Default                |
| ------------------ | ---------------------------------------------- | ---------------------- |
| `DB_URL`           | PostgreSQL connection string                   | Required               |
| `JWT_SECRET`       | JWT signing secret (min 32 chars)              | Required               |
| `JWT_EXPIRY_HOURS` | Token expiration time                          | 24                     |
| `PORT`             | Server port                                    | 8080                   |
| `LOG_LEVEL`        | Logging level (debug, info, warn, error)       | info                   |
| `ENVIRONMENT`      | Environment (development, staging, production) | development            |
| `REDIS_URL`        | Redis connection string                        | redis://localhost:6379 |
| `NATS_URL`         | NATS connection string                         | nats://localhost:4222  |
| `RATE_LIMIT_RPS`   | Requests per second limit                      | 10                     |
| `ALLOWED_ORIGINS`  | CORS allowed origins                           | \*                     |

## Development

### Project Structure

The project follows clean architecture principles:

- **Domain Layer**: Business entities and errors
- **Repository Layer**: Data access with interfaces
- **Service Layer**: Business logic
- **Handler Layer**: HTTP request handling
- **Infrastructure**: External services (cache, messaging)

### Key Design Patterns

1. **Dependency Injection**: All dependencies injected via constructors
2. **Interface Segregation**: Small, focused interfaces
3. **Repository Pattern**: Data access abstraction
4. **Service Layer**: Business logic separation
5. **Middleware Chain**: Composable HTTP middleware
6. **Error Handling**: Domain-specific errors with HTTP mapping

### Running Tests

```bash
# Run all tests
make test

# Run integration tests
make test-integration

# View coverage
make test
open coverage.html
```

### Code Quality

```bash
# Format code
make fmt

# Run linter
make lint

# Tidy dependencies
make tidy
```

## Database Migrations

### Create Migration

```bash
make migrate-create name=add_email_verification
```

### Run Migrations

```bash
make migrate-up
```

### Rollback Migrations

```bash
make migrate-down
```

## Monitoring & Observability

### Prometheus Metrics

Available at `/metrics`:

- HTTP request duration histograms
- HTTP request counters (by path, method, status)
- Database query duration
- Database query counters (by operation, status)

### Structured Logging

All logs are structured JSON (in production) or console (in development):

```json
{
  "level": "info",
  "time": "2025-01-01T12:00:00Z",
  "message": "request completed",
  "method": "POST",
  "path": "/api/v1/login",
  "status": 200,
  "duration": 45.2
}
```

### Health Checks

- `/health` - Checks all dependencies (DB, Redis, NATS)
- `/ready` - Kubernetes readiness probe
- `/live` - Kubernetes liveness probe

## Production Deployment

### Docker

```bash
# Build image
docker build -t user-auth-app:latest .

# Run container
docker run -p 8080:8080 \
  -e DB_URL="..." \
  -e JWT_SECRET="..." \
  user-auth-app:latest
```

### Kubernetes

See `k8s/` directory for Kubernetes manifests (deployment, service, ingress, configmap).

### Environment-Specific Configuration

1. **Development**: Uses console logging, relaxed CORS
2. **Staging**: JSON logging, stricter rate limits
3. **Production**: JSON logging, strict CORS, enhanced security

## Security Considerations

- ✅ Passwords hashed with bcrypt
- ✅ JWT tokens with expiration
- ✅ Rate limiting per IP
- ✅ CORS protection
- ✅ Input validation
- ✅ SQL injection protection (via sqlc)
- ✅ Panic recovery middleware
- ✅ Secure headers

## Performance Features

- Redis caching with fallback to in-memory
- Connection pooling for PostgreSQL
- Efficient database queries via sqlc
- Request timeout handling
- Graceful shutdown
- Prometheus metrics for monitoring

## Troubleshooting

### Database Connection Issues

```bash
# Check database is running
docker ps | grep postgres

# Test connection
psql -h localhost -U postgres -d dbname
```

### Redis Issues

```bash
# Check Redis
docker ps | grep redis

# Test connection
redis-cli -h localhost ping
```

### Port Already in Use

```bash
# Find process using port 8080
lsof -i :8080

# Kill process
kill -9 <PID>
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run `make lint` and `make test`
6. Submit a pull request

## License

MIT License - see LICENSE file for details

## Support

For issues and questions:

- Create an issue on GitHub
- Check existing documentation
- Review API examples

## Roadmap

- [ ] OAuth2 integration
- [ ] Email verification
- [ ] Password reset flow
- [ ] 2FA support
- [ ] API rate limiting per user
- [ ] Request ID tracing
- [ ] OpenAPI/Swagger documentation
- [ ] GraphQL support
- [ ] Websocket support
