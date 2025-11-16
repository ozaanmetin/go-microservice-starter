package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ozaanmetin/go-microservice-starter/internal/api"
	"github.com/ozaanmetin/go-microservice-starter/internal/api/middlewares"
	"github.com/ozaanmetin/go-microservice-starter/internal/config"
	"github.com/ozaanmetin/go-microservice-starter/pkg/logging"
)

func main() {
	// Setup configuration
	cfg, err := config.Load()

	if err != nil {
		panic(err)
	}

	// Setup logger
	_, err = logging.Init(logging.Config{
		Level:  cfg.Logger.Level,
		Format: cfg.Logger.Format,
	})

	if err != nil {
		panic(err)
	}
	defer logging.L().Sync()

	// Create Fiber app with custom error handler
	app := fiber.New(
		fiber.Config{
			AppName:      cfg.Server.AppName,
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
			ErrorHandler: middlewares.ServiceErrorErrorHandler(),
		},
	)

	// Setup middlewares
	api.SetupMiddlewares(app, cfg)
	// Setup routes
	api.SetupRoutes(app)

	app.Get("/slow", func(c *fiber.Ctx) error {
		time.Sleep(5 * time.Second) // 5 saniye bekle
		return c.JSON(fiber.Map{"message": "slow response"})
	})

	// Start server in goroutine
	go func() {
		logging.L().WithField("host", cfg.Server.Host).WithField("port", cfg.Server.Port).Info("Starting server...")
		if err := app.Listen(config.GetServerAddress(cfg)); err != nil {
			logging.L().WithError(err).Fatal("Failed to start server")
		}
	}()

	// Handle graceful shutdown (blocks here)
	gracefulShutdown(app, cfg.Server.ShutdownTimeout)
}

func gracefulShutdown(app *fiber.App, shutdownTimeout time.Duration) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	logging.L().Info("Shutting down server...")

	if err := app.ShutdownWithTimeout(shutdownTimeout); err != nil {
		logging.L().WithError(err).Error("Error during server shutdown")
	}

	logging.L().Info("Server gracefully stopped!")
}
