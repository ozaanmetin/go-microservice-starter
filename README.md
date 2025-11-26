# Go Microservice Starter

A production-ready Go microservice starter template built with clean architecture principles and modern best practices.

## Features

### Core Architecture
- **Clean Architecture**: Separation of HTTP layer, business logic, and infrastructure
- **Generic Handler Pattern**: Framework-agnostic business logic with automatic request parsing
- **Centralized Error Handling**: Custom `ServiceError` type with consistent error responses
- **Graceful Shutdown**: Configurable timeout with proper signal handling

### Observability
- **Prometheus Metrics**: HTTP requests, duration, and in-flight metrics
- **Structured Logging**: Zap-based logging with JSON/console output
- **Request Tracing**: Unique request IDs for distributed tracing

### Resilience
- **Circuit Breaker**: Protection against cascading failures using `sony/gobreaker`
- **Rate Limiting**: Global and per-endpoint rate limiting with Redis backend
- **Panic Recovery**: Automatic panic recovery middleware

### Infrastructure
- **Redis Integration**: Client and storage implementations for caching and rate limiting
- **Docker Compose**: Complete stack with Prometheus, and Redis
- **Configuration Management**: Viper-based config with environment variable support

## Quick Start

```bash
# Run locally
make runserver

# Or with Docker
docker-compose up -d
```

**Endpoints**:
- App: http://localhost:8000
- Health Check: http://localhost:8000/healthcheck
- Metrics: http://localhost:8000/metrics
- Prometheus: http://localhost:9090


## Future Enhancements

- [ ] Database integration
- [ ] Authentication & Authorization (JWT, OAuth2)
- [ ] Swagger/OpenAPI documentation
- [ ] Migration system
- [ ] CI/CD pipeline configuration
- [ ] Kubernetes manifests
