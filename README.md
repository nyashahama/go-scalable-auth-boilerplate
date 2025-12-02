# User Authentication Service

A production-ready Go authentication service template with JWT authentication, PostgreSQL, Redis caching, and NATS messaging. Built with clean architecture principles and modern Go best practices.

## 🚀 Features

- **JWT Authentication** - Secure token-based authentication with configurable expiry
- **User Management** - Registration, login, and profile endpoints
- **PostgreSQL Database** - Type-safe database queries with sqlc
- **Redis Caching** - Fast profile lookups with automatic fallback to in-memory cache
- **NATS Messaging** - Async event publishing (e.g., email verification)
- **Prometheus Metrics** - Built-in monitoring and observability
- **Rate Limiting** - IP-based rate limiting to prevent abuse
- **CORS Support** - Configurable cross-origin resource sharing
- **Input Validation** - Comprehensive request validation
- **Structured Logging** - JSON logging with zerolog
- **Docker Support** - Full containerization with docker-compose
- **Graceful Shutdown** - Clean service termination
- **Health Checks** - Ready-to-use health and metrics endpoints

## 📋 Prerequisites

- Go 1.23 or higher
- Docker & Docker Compose
- PostgreSQL 16 (or use Docker)
- Redis 7 (or use Docker)
- NATS 2 (or use Docker)

## 🛠️ Installation

### 1. Clone the Repository

```bash
git clone https://github.com/nyashahama/go-scalable-auth-boilerplate.git
cd user-auth-app
```

### 2. Install Dependencies

```bash
# Install Go dependencies
go mod download

# Install development tools
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### 3. Setup Environment

```bash
# Generate .env file with secure JWT secret
chmod +x setup-env.sh
./setup-env.sh

# Or manually create .env (see .env.example)
cp .env.example .env
# Edit .env and set your JWT_SECRET
```

### 4. Start Services

```bash
# Option A: Use the startup script (recommended)
chmod +x start.sh
./start.sh

# Option B: Manual setup
docker-compose up -d postgres redis nats
make migrate-up
make sqlc
```

### 5. Run the Application

```bash
# Development mode
make run

# Or build and run
make build
./bin/api

# Or with Docker
docker-compose up -d
```

The server will start on `http://localhost:8080`

## 🔧 Configuration

All configuration is done via environment variables in `.env`:

| Variable           | Description                              | Default               |
| ------------------ | ---------------------------------------- | --------------------- |
| `DB_URL`           | PostgreSQL connection string             | Required              |
| `JWT_SECRET`       | Secret key for JWT signing               | Required              |
| `JWT_EXPIRY_HOURS` | JWT token expiry time                    | 24                    |
| `PORT`             | Server port                              | 8080                  |
| `LOG_LEVEL`        | Logging level (debug, info, warn, error) | info                  |
| `ENVIRONMENT`      | Environment (development, production)    | development           |
| `TIMEOUT_SECONDS`  | Request timeout                          | 30                    |
| `REDIS_URL`        | Redis connection string                  | localhost:6379        |
| `NATS_URL`         | NATS connection string                   | nats://localhost:4222 |
| `ALLOWED_ORIGINS`  | CORS allowed origins (comma-separated)   | \*                    |
| `RATE_LIMIT_RPS`   | Rate limit requests per second           | 10                    |
| `RATE_LIMIT_BURST` | Rate limit burst size                    | 20                    |

### Generate Secure JWT Secret

```bash
# Generate a 32-byte base64 encoded secret
openssl rand -base64 32

# Or use the makefile
make generate-jwt
```

## 📚 API Documentation

### Public Endpoints

#### Health Check

```bash
GET /health
```

#### Register User

```bash
POST /register
Content-Type: application/json

{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "secure_password123",
  "role": "user"
}
```

**Response:**

```json
{
  "id": 1,
  "username": "johndoe",
  "email": "john@example.com",
  "role": "user",
  "created_at": "2025-12-02T10:30:00Z"
}
```

#### Login

```bash
POST /login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "secure_password123"
}
```

**Response:**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Protected Endpoints

#### Get User Profile

```bash
GET /users/{id}
Authorization: Bearer {token}
```

**Response:**

```json
{
  "id": 1,
  "username": "johndoe",
  "email": "john@example.com",
  "role": "user",
  "created_at": "2025-12-02T10:30:00Z"
}
```

### Monitoring Endpoints

#### Prometheus Metrics

```bash
GET /metrics
```

## 🧪 Testing

### Run Tests

```bash
# Run all tests
make test

# Run with coverage
make test

# Run integration tests
make test-integration

# View coverage report
open coverage.html
```

### Manual API Testing

Use the provided test script:

```bash
chmod +x test-api.sh
./test-api.sh
```

Or use curl:

```bash
# Register
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"password123","role":"user"}'

# Login
TOKEN=$(curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}' | jq -r '.token')

# Get profile
curl http://localhost:8080/users/1 \
  -H "Authorization: Bearer $TOKEN"
```

## 🗂️ Project Structure

```
user-auth-app/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/                  # Configuration management
│   ├── domain/                  # Domain models
│   ├── handlers/                # HTTP handlers
│   ├── middleware/              # HTTP middleware (auth, CORS, rate limit)
│   ├── repository/              # Data access layer
│   │   ├── sqlc/               # Generated database code
│   │   ├── schema.sql          # Database schema
│   │   └── queries.sql         # SQL queries
│   ├── services/                # Business logic
│   └── validator/               # Input validation
├── migrations/                  # Database migrations
├── scripts/                     # Utility scripts
├── .env.example                 # Example environment config
├── docker-compose.yml           # Docker services configuration
├── Dockerfile                   # Application container
├── Makefile                     # Common tasks
├── sqlc.yaml                    # sqlc configuration
└── README.md
```

## 🐳 Docker Usage

### Using Docker Compose

```bash
# Start all services (Postgres, Redis, NATS, API)
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Rebuild API container
docker-compose build api
docker-compose up -d api
```

### Individual Services

```bash
# Start only infrastructure services
docker-compose up -d postgres redis nats

# Then run the API locally for development
make run
```

## 📊 Database Migrations

### Create Migration

```bash
make migrate-create name=add_users_table
```

### Run Migrations

```bash
# Migrate up
make migrate-up

# Migrate down
make migrate-down
```

### Generate sqlc Code

```bash
make sqlc
```

## 🔍 Common Tasks

```bash
# Format code
make fmt

# Run linter
make lint

# Tidy dependencies
make tidy

# Clean build artifacts
make clean

# Build binary
make build

# Run application
make run

# Generate JWT secret
make generate-jwt
```

## 🚀 Production Deployment

### Environment Variables

Ensure these are set in production:

- `ENVIRONMENT=production`
- `LOG_LEVEL=info`
- Strong `JWT_SECRET` (32+ bytes)
- Proper `ALLOWED_ORIGINS`
- Database connection with SSL: `?sslmode=require`

### Docker Deployment

```bash
# Build for production
docker build -t user-auth-app:latest .

# Run with production config
docker run -d \
  --name user-auth-app \
  -p 8080:8080 \
  --env-file .env.production \
  user-auth-app:latest
```

### Kubernetes

Coming soon! Check the `k8s/` directory for Kubernetes manifests.

## 🔒 Security Best Practices

- ✅ JWT tokens with expiration
- ✅ Password hashing with bcrypt
- ✅ Rate limiting on all endpoints
- ✅ Input validation and sanitization
- ✅ CORS configuration
- ✅ SQL injection prevention (sqlc)
- ✅ Secure headers
- ✅ Environment-based secrets

**⚠️ Production Checklist:**

- [ ] Use strong JWT secret (32+ bytes)
- [ ] Enable HTTPS/TLS
- [ ] Configure proper CORS origins
- [ ] Set up database backups
- [ ] Monitor logs and metrics
- [ ] Implement request logging
- [ ] Add rate limiting per user
- [ ] Set up alerting

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [sqlc](https://github.com/sqlc-dev/sqlc) - Type-safe SQL in Go
- [chi](https://github.com/go-chi/chi) - Lightweight HTTP router
- [zerolog](https://github.com/rs/zerolog) - Zero-allocation JSON logger
- [pgx](https://github.com/jackc/pgx) - PostgreSQL driver
- [jwt-go](https://github.com/golang-jwt/jwt) - JWT implementation

## 📧 Support

For issues and questions:

- 🐛 [Report bugs](https://github.com/nyashahama/go-scalable-auth-boilerplate/issues)
- 💡 [Request features](https://github.com/nyashahama/go-scalable-auth-boilerplate/issues)
- 📖 [Documentation](https://github.com/nyashahama/user-auth-app/wiki)

## 🗺️ Roadmap

- [ ] Refresh token support
- [ ] Email verification implementation
- [ ] Password reset flow
- [ ] OAuth2 integration (Google, GitHub)
- [ ] Role-based access control (RBAC)
- [ ] Two-factor authentication (2FA)
- [ ] API versioning
- [ ] GraphQL support
- [ ] WebSocket support
- [ ] Admin dashboard

---
