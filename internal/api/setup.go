package api

import (
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/ozaanmetin/go-microservice-starter/internal/api/healthcheck"
	"github.com/ozaanmetin/go-microservice-starter/internal/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	// Middleware imports
	middlewares "github.com/ozaanmetin/go-microservice-starter/internal/api/middlewares"

	// Infrastructure imports
	infrahttp "github.com/ozaanmetin/go-microservice-starter/internal/infrastructure/http"
)

// SetupRoutes configures all HTTP routes for the application
func SetupRoutes(app *fiber.App) {
	// Handlers
	healthHandler := healthcheck.NewHealthCheckHandler()

	// Setup routes
	app.Get("/healthcheck", infrahttp.AdaptHandler(healthHandler))

	// Prometheus metrics endpoint
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
}

func SetupMiddlewares(app *fiber.App, cfg *config.Config) {
	// Middleware order matters!

	app.Use(middlewares.RequestID())
	app.Use(middlewares.Recover())
	app.Use(middlewares.Metrics())
	app.Use(middlewares.Logger(middlewares.LoggerConfig{
		SkipPaths: []string{"/metrics"}, // Skip metrics endpoint to reduce noise
	}))
}
