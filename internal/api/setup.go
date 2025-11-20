package api

import (
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/ozaanmetin/go-microservice-starter/internal/api/circuit_breaker_example"
	"github.com/ozaanmetin/go-microservice-starter/internal/api/healthcheck"
	"github.com/ozaanmetin/go-microservice-starter/internal/config"
	infraredis "github.com/ozaanmetin/go-microservice-starter/internal/infrastructure/redis"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	// Middleware imports
	middlewares "github.com/ozaanmetin/go-microservice-starter/internal/api/middlewares"

	// Infrastructure imports
	infrahttp "github.com/ozaanmetin/go-microservice-starter/internal/infrastructure/http"
)

// SetupRoutes configures all HTTP routes for the application
func SetupRoutes(app *fiber.App, cfg *config.Config) {
	// Handlers
	healthHandler := healthcheck.NewHealthCheckHandler()
	circuitBreakerExampleHandler := circuitBreakerExample.NewExampleHandler()

	// Rate Limiter
	// TODO: Only for demonstration, adjust as needed or remove from healthcheck
	healthCheckRateLimiter := middlewares.NewEndpointRateLimiter(
		10,                                     // Max too many requests
		1*time.Minute,                          // Per minute
		infraredis.NewStorage(&cfg.Redis),      // Create storage here
		middlewares.KeyByIP,
	)

	// Setup routes with endpoint-specific rate limiter
	app.Get("/healthcheck", healthCheckRateLimiter, infrahttp.AdaptHandler(healthHandler))

	// Example endpoint with circuit breaker
	app.Get("/circuit-breaker-example", infrahttp.AdaptHandler(circuitBreakerExampleHandler))

	// Prometheus metrics endpoint
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
}

func SetupMiddlewares(app *fiber.App, cfg *config.Config) {
	// Middleware order matters!

	app.Use(middlewares.RequestID())
	app.Use(middlewares.Recover())

	// Setup rate limiter if enabled
	if cfg.Server.RateLimiter.Enabled {
		app.Use(middlewares.RateLimiter(middlewares.RateLimiterConfig{
			Max:          cfg.Server.RateLimiter.Max,
			Expiration:   cfg.Server.RateLimiter.Expiration,
			Storage:      infraredis.NewStorage(&cfg.Redis),
			KeyGenerator: middlewares.KeyByIP,
			SkipPaths:    []string{},
		}))
	}

	app.Use(middlewares.Metrics())
	app.Use(middlewares.Logger(middlewares.LoggerConfig{
		SkipPaths: []string{"/metrics"}, // Skip metrics endpoint to reduce noise
	}))
}
