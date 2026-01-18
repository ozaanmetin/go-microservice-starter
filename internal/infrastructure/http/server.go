package http

import (
	"net/http"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"

	"github.com/ozaanmetin/go-microservice-starter/internal/config"
	"github.com/ozaanmetin/go-microservice-starter/internal/infrastructure/http/middlewares"
	"github.com/ozaanmetin/go-microservice-starter/internal/infrastructure/redis"
	"github.com/ozaanmetin/go-microservice-starter/pkg/logging"
)

// RouteSetupFunc is injected by the api layer to register routes.
type RouteSetupFunc func(s *Server)


// Structs in order to abstract the fiber.App and fiber.Router
type Server struct {
	app *fiber.App
	cfg *config.Config
}

type RouteGroup struct {
	router fiber.Router
}

func NewServer(cfg *config.Config, setupRoutes RouteSetupFunc) *Server {
	app := fiber.New(fiber.Config{
		AppName:      cfg.Server.AppName,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		ErrorHandler: middlewares.ServiceErrorErrorHandler(),
	})

	server := &Server{
		app: app,
		cfg: cfg,
	}

	server.setupMiddlewares()
	setupRoutes(server)

	return server
}

func (s *Server) Listen(address string) error {
	logging.L().WithField("server_address", address).Info("Starting server...")
	return s.app.Listen(address)
}

func (s *Server) Shutdown(timeout time.Duration) error {
	return s.app.ShutdownWithTimeout(timeout)
}

// ---- Server Routes

func (s *Server) Get(path string, handler fiber.Handler, mws ...fiber.Handler) {
	handlers := append(mws, handler)
	s.app.Get(path, handlers...)
}

func (s *Server) Post(path string, handler fiber.Handler, mws ...fiber.Handler) {
	handlers := append(mws, handler)
	s.app.Post(path, handlers...)
}

func (s *Server) Put(path string, handler fiber.Handler, mws ...fiber.Handler) {
	handlers := append(mws, handler)
	s.app.Put(path, handlers...)
}

func (s *Server) Delete(path string, handler fiber.Handler, mws ...fiber.Handler) {
	handlers := append(mws, handler)
	s.app.Delete(path, handlers...)
}

func (s *Server) Group(prefix string, mws ...fiber.Handler) *RouteGroup {
	return &RouteGroup{router: s.app.Group(prefix, mws...)}
}

// Mount adapts standard http.Handler to all HTTP methods
func (s *Server) Mount(path string, handler http.Handler) {
	s.app.All(path, adaptor.HTTPHandler(handler))
}

// ---- RouteGroup Routes

func (g *RouteGroup) Get(path string, handler fiber.Handler, mws ...fiber.Handler) {
	handlers := append(mws, handler)
	g.router.Get(path, handlers...)
}

func (g *RouteGroup) Post(path string, handler fiber.Handler, mws ...fiber.Handler) {
	handlers := append(mws, handler)
	g.router.Post(path, handlers...)
}

func (g *RouteGroup) Put(path string, handler fiber.Handler, mws ...fiber.Handler) {
	handlers := append(mws, handler)
	g.router.Put(path, handlers...)
}

func (g *RouteGroup) Delete(path string, handler fiber.Handler, mws ...fiber.Handler) {
	handlers := append(mws, handler)
	g.router.Delete(path, handlers...)
}

// ---- Middlewares

func (s *Server) setupMiddlewares() {
	s.app.Use(middlewares.RequestID())
	s.app.Use(middlewares.Recover())

	if s.cfg.Server.RateLimiter.Enabled {
		s.app.Use(middlewares.RateLimiter(middlewares.RateLimiterConfig{
			Max:          s.cfg.Server.RateLimiter.Max,
			Expiration:   s.cfg.Server.RateLimiter.Expiration,
			Storage:      redis.NewStorage(&s.cfg.Redis),
			KeyGenerator: middlewares.KeyByIP,
			SkipPaths:    []string{},
		}))
	}

	s.app.Use(middlewares.Metrics())
	s.app.Use(middlewares.Logger(middlewares.LoggerConfig{
		SkipPaths: []string{"/metrics"},
	}))
}
