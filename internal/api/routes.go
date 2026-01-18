package api

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ozaanmetin/go-microservice-starter/internal/api/features/auth"
	"github.com/ozaanmetin/go-microservice-starter/internal/api/features/circuit_breaker_example"
	"github.com/ozaanmetin/go-microservice-starter/internal/api/features/healthcheck"
	"github.com/ozaanmetin/go-microservice-starter/internal/api/features/profile"
	"github.com/ozaanmetin/go-microservice-starter/internal/config"
	"github.com/ozaanmetin/go-microservice-starter/internal/domain/user"
	"github.com/ozaanmetin/go-microservice-starter/internal/infrastructure/http/middlewares"
	infrahttp "github.com/ozaanmetin/go-microservice-starter/internal/infrastructure/http"
	infraredis "github.com/ozaanmetin/go-microservice-starter/internal/infrastructure/redis"
	pkgJWT "github.com/ozaanmetin/go-microservice-starter/pkg/jwt"
)

// NewRouteSetup creates a route setup function with the given dependencies
// This returns a function that can be passed to infrahttp.NewServer
func NewRouteSetup(cfg *config.Config, db *sqlx.DB) infrahttp.RouteSetupFunc {
	return func(s *infrahttp.Server) {
		// Initialize JWT Manager
		jwtManager := pkgJWT.NewManager(
			cfg.JWT.Secret,
			cfg.JWT.AccessTokenDuration,
			cfg.JWT.RefreshTokenDuration,
		)

		// Initialize repositories
		userRepo := user.NewRepository(db)

		// Initialize services
		authService := auth.NewAuthService(userRepo, jwtManager)
		profileService := profile.NewProfileService(userRepo)

		// Initialize handlers
		refreshTokenHandler := auth.NewRefreshTokenHandler(authService)
		loginHandler := auth.NewLoginHandler(authService)
		registerHandler := auth.NewRegisterHandler(authService)
		profileHandler := profile.NewGetProfileHandler(profileService)
		healthHandler := healthcheck.NewHealthCheckHandler()
		circuitBreakerExampleHandler := circuitBreakerExample.NewExampleHandler()

		// Rate Limiter for healthcheck
		healthCheckRateLimiter := middlewares.NewEndpointRateLimiter(
			10,
			1*time.Minute,
			infraredis.NewStorage(&cfg.Redis),
			middlewares.KeyByIP,
		)

		// Public routes
		s.Get("/healthcheck", infrahttp.AdaptHandler(healthHandler), healthCheckRateLimiter)
		s.Get("/circuit-breaker-example", infrahttp.AdaptHandler(circuitBreakerExampleHandler))
		s.Mount("/metrics", promhttp.Handler())

		// Auth routes (public)
		authGroup := s.Group("/auth")
		authGroup.Post("/register", infrahttp.AdaptHandler(registerHandler))
		authGroup.Post("/login", infrahttp.AdaptHandler(loginHandler))
		authGroup.Post("/refresh", infrahttp.AdaptHandler(refreshTokenHandler))

		// Protected routes (require JWT authentication)
		apiGroup := s.Group("/api", middlewares.AuthMiddleware(jwtManager))
		apiGroup.Get("/profile", infrahttp.AdaptHandler(profileHandler))
	}
}
