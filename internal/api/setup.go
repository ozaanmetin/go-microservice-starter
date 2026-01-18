package api

import (
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/ozaanmetin/go-microservice-starter/internal/api/auth"
	"github.com/ozaanmetin/go-microservice-starter/internal/api/circuit_breaker_example"
	"github.com/ozaanmetin/go-microservice-starter/internal/api/healthcheck"
	"github.com/ozaanmetin/go-microservice-starter/internal/api/profile"
	"github.com/ozaanmetin/go-microservice-starter/internal/config"
	"github.com/ozaanmetin/go-microservice-starter/internal/domain/user"
	infraredis "github.com/ozaanmetin/go-microservice-starter/internal/infrastructure/redis"
	pkgJWT "github.com/ozaanmetin/go-microservice-starter/pkg/jwt"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	// Middleware imports
	"github.com/ozaanmetin/go-microservice-starter/internal/middlewares"

	// Infrastructure imports
	infrahttp "github.com/ozaanmetin/go-microservice-starter/internal/infrastructure/http"
)

// SetupRoutes configures all HTTP routes for the application
func SetupRoutes(app *fiber.App, cfg *config.Config, db *sqlx.DB) {
	// Initialize JWT Manager
	jwtManager := pkgJWT.NewManager(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenDuration,
		cfg.JWT.RefreshTokenDuration,
	)

	// Initialize repositories
	// -------------------------------------------------------------
	// User
	userRepo := user.NewRepository(db)

	// Initialize services
	// -------------------------------------------------------------
	// Auth
	authService := auth.NewAuthService(userRepo, jwtManager)
	// Profile
	profileService := profile.NewProfileService(userRepo)

	// Initialize handlers
	// -------------------------------------------------------------
	// Login
	refreshTokenHandler := auth.NewRefreshTokenHandler(authService)
	loginHandler := auth.NewLoginHandler(authService)
	// Register
	registerHandler := auth.NewRegisterHandler(authService)
	// Profile
	profileHandler := profile.NewGetProfileHandler(profileService)


	// Other handlers
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

	// Public routes
	app.Get("/healthcheck", healthCheckRateLimiter, infrahttp.AdaptHandler(healthHandler))
	app.Get("/circuit-breaker-example", infrahttp.AdaptHandler(circuitBreakerExampleHandler))
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	// Auth routes (public)
	authGroup := app.Group("/auth")
	authGroup.Post("/register", infrahttp.AdaptHandler(registerHandler))
	authGroup.Post("/login", infrahttp.AdaptHandler(loginHandler))
	authGroup.Post("/refresh", infrahttp.AdaptHandler(refreshTokenHandler))

	// Protected routes (require JWT authentication)
	api := app.Group("/api", middlewares.AuthMiddleware(jwtManager))
	// Profile routes
	api.Get("/profile", infrahttp.AdaptHandler(profileHandler))
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
