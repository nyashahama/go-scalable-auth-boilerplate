# User Authentication API

A production-ready authentication API built with Go, featuring JWT authentication, Redis caching, NATS messaging, email notifications, and comprehensive observability.

## Features

- 🔐 **JWT Authentication** - Secure token-based authentication
- 👤 **User Management** - Registration, login, profile management
- 📧 **Email Notifications** - AWS SES integration with beautiful HTML templates
- 🚀 **High Performance** - Redis caching for optimized data access
- 📨 **Event-Driven** - NATS messaging for async operations
- 📊 **Observability** - Prometheus metrics, structured logging
- 🛡️ **Security** - Rate limiting, CORS, input validation
- 🗃️ **Clean Architecture** - Separation of concerns, dependency injection
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
├── email/            # Email service (SES/SMTP)
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
- AWS account (optional, for production email)

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

# Access Mailpit (local email UI)
open http://localhost:8025
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

### 7. Test Email (Local Development)

```bash
# Register a user
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "SecurePass123!",
    "role": "user"
  }'

# Check email at http://localhost:8025
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

# Response: 201 Created
# Triggers: Welcome email sent automatically
```

#### Login

```bash
POST /api/v1/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "securepassword123"
}

# Response: 200 OK with JWT token
# Triggers: Login alert email (optional security feature)
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
GET /health      # Comprehensive health check (includes email service)
GET /ready       # Readiness probe
GET /live        # Liveness probe
GET /metrics     # Prometheus metrics
```

## Email Service

The application includes a robust email service that works in both development and production.

### 📧 Supported Email Types

1. **Welcome Email** - Sent automatically on user registration
2. **Password Reset** - For password recovery (ready to implement)
3. **Email Verification** - For account verification (ready to implement)
4. **Password Changed** - Security notification
5. **Login Alert** - Suspicious login notifications

### 🏠 Local Development (SMTP with Mailpit)

Emails are caught by Mailpit - nothing sent to real addresses:

- **SMTP Server:** `localhost:1025`
- **Web UI:** http://localhost:8025
- **Configuration:** Already set in `docker-compose.yml`

```bash
# Start services
docker-compose up -d

# Test email
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"Pass123!"}'

# View email in browser
open http://localhost:8025
```

### 🚀 Production (AWS SES)

For production deployment (Render, AWS, etc.), use AWS SES:

**Quick Setup (15 minutes):**

1. **Verify email in AWS SES Console:**
   - Go to https://console.aws.amazon.com/ses/
   - Verified identities → Create identity → Email address
   - Click verification link in email

2. **Create IAM user:**
   - IAM Console → Users → Create user
   - Attach policy: `AmazonSESFullAccess`
   - Create access key → Save credentials

3. **Configure environment variables:**
   ```bash
   EMAIL_PROVIDER=ses
   AWS_REGION=us-east-1
   AWS_ACCESS_KEY_ID=AKIA...
   AWS_SECRET_ACCESS_KEY=wJal...
   EMAIL_FROM_ADDRESS=noreply@yourdomain.com
   EMAIL_FROM_NAME=Your App
   ```

4. **Deploy and test!**

📚 **Detailed Setup Guide:** See [EMAIL_SETUP.md](./EMAIL_SETUP.md) for complete instructions.

🚀 **Quick Start Guide:** See [QUICKSTART_EMAIL.md](./QUICKSTART_EMAIL.md) for 5-minute setup.

### Email Configuration Options

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `EMAIL_PROVIDER` | Email provider (`ses` or `smtp`) | `smtp` | Yes |
| `EMAIL_FROM_ADDRESS` | Sender email address | `noreply@localhost` | Yes |
| `EMAIL_FROM_NAME` | Sender display name | `Auth Service` | Yes |
| `AWS_REGION` | AWS region for SES | `us-east-1` | For SES |
| `AWS_ACCESS_KEY_ID` | AWS access key | - | For SES |
| `AWS_SECRET_ACCESS_KEY` | AWS secret key | - | For SES |
| `SMTP_HOST` | SMTP server hostname | `localhost` | For SMTP |
| `SMTP_PORT` | SMTP server port | `1025` | For SMTP |
| `SMTP_USERNAME` | SMTP username (optional) | - | No |
| `SMTP_PASSWORD` | SMTP password (optional) | - | No |

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
- **Infrastructure**: External services (cache, messaging, email)

### Key Design Patterns

1. **Dependency Injection**: All dependencies injected via constructors
2. **Interface Segregation**: Small, focused interfaces
3. **Repository Pattern**: Data access abstraction
4. **Service Layer**: Business logic separation
5. **Middleware Chain**: Composable HTTP middleware
6. **Error Handling**: Domain-specific errors with HTTP mapping
7. **Graceful Degradation**: Optional services don't block core functionality

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

- `/health` - Checks all dependencies (DB, Redis, NATS, Email)
- `/ready` - Kubernetes readiness probe
- `/live` - Kubernetes liveness probe

Example health check response:

```json
{
  "status": "healthy",
  "timestamp": "2025-01-01T12:00:00Z",
  "services": {
    "database": "healthy",
    "cache": "healthy",
    "messaging": "healthy",
    "email": "healthy"
  },
  "version": "1.0.0"
}
```

## Production Deployment

### Docker

```bash
# Build image
docker build -t user-auth-app:latest .

# Run container
docker run -p 8080:8080 \
  -e DB_URL="..." \
  -e JWT_SECRET="..." \
  -e EMAIL_PROVIDER="ses" \
  -e AWS_ACCESS_KEY_ID="..." \
  -e AWS_SECRET_ACCESS_KEY="..." \
  user-auth-app:latest
```

### Render

1. Create new Web Service
2. Connect your repository
3. Add environment variables (see `.env.example`)
4. For email, use AWS SES (SMTP ports are blocked on Render)

### Kubernetes

See `k8s/` directory for Kubernetes manifests (deployment, service, ingress, configmap).

### Environment-Specific Configuration

1. **Development**: 
   - Console logging
   - Relaxed CORS
   - SMTP email with Mailpit
   
2. **Staging**: 
   - JSON logging
   - Stricter rate limits
   - AWS SES email
   
3. **Production**: 
   - JSON logging
   - Strict CORS
   - Enhanced security
   - AWS SES email
   - Domain-verified emails

## Security Considerations

- ✅ Passwords hashed with bcrypt
- ✅ JWT tokens with expiration
- ✅ Rate limiting per IP
- ✅ CORS protection
- ✅ Input validation
- ✅ SQL injection protection (via sqlc)
- ✅ Panic recovery middleware
- ✅ Secure headers
- ✅ Email sent asynchronously (non-blocking)
- ✅ Graceful degradation (email failures don't block auth)

## Performance Features

- Redis caching with fallback to in-memory
- Connection pooling for PostgreSQL
- Efficient database queries via sqlc
- Request timeout handling
- Graceful shutdown
- Prometheus metrics for monitoring
- Async email sending (non-blocking)
- AWS SES for reliable email delivery

## Email Costs (AWS SES)

- **Free Tier**: 62,000 emails/month (forever)
- **After Free Tier**: $0.10 per 1,000 emails

**Example costs:**
- 10,000 emails/month: ~$0.70
- 100,000 emails/month: ~$9.70

Much cheaper than SendGrid, Mailgun, or Postmark!

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

### Email Not Sending

```bash
# Check health endpoint
curl http://localhost:8080/health

# Check logs for email service status
docker logs auth_api | grep email

# For local dev, check Mailpit UI
open http://localhost:8025
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
- See [EMAIL_SETUP.md](./EMAIL_SETUP.md) for email configuration

## Documentation

- [README.md](./README.md) - This file
- [EMAIL_SETUP.md](./EMAIL_SETUP.md) - Complete email setup guide
- [QUICKSTART_EMAIL.md](./QUICKSTART_EMAIL.md) - 5-minute email quick start
- [.env.example](./.env.example) - Environment variables reference

## Roadmap

- [x] JWT authentication
- [x] Redis caching
- [x] NATS messaging
- [x] Email notifications (SES/SMTP)
- [x] Prometheus metrics
- [x] Health checks
- [ ] OAuth2 integration (Google, GitHub)
- [ ] Email verification flow
- [ ] Password reset flow (email integration ready)
- [ ] 2FA support
- [ ] API rate limiting per user
- [ ] Request ID tracing
- [ ] OpenAPI/Swagger documentation
- [ ] GraphQL support
- [ ] Websocket support
- [ ] Admin dashboard

## Tech Stack

- **Language:** Go 1.21+
- **Web Framework:** Chi Router
- **Database:** PostgreSQL 16 with sqlc
- **Cache:** Redis 7
- **Messaging:** NATS 2 with JetStream
- **Email:** AWS SES / SMTP
- **Auth:** JWT (golang-jwt)
- **Logging:** Zerolog (structured logging)
- **Metrics:** Prometheus
- **Validation:** Custom validator
- **Testing:** Go testing + testify
- **Containerization:** Docker & Docker Compose
- **Orchestration:** Kubernetes ready

## Performance Benchmarks

Coming soon - including:
- Requests per second
- Average latency
- Database query performance
- Cache hit rates
- Email delivery times

---

**Built with ❤️ using Go** | [Report Bug](https://github.com/nyashahama/go-scalable-auth-boilerplate/issues) | [Request Feature](https://github.com/nyashahama/go-scalable-auth-boilerplate/issues)
