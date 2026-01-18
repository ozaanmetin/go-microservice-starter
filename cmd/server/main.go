package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ozaanmetin/go-microservice-starter/pkg/logging"
	"github.com/ozaanmetin/go-microservice-starter/internal/api"
	"github.com/ozaanmetin/go-microservice-starter/internal/config"
	"github.com/ozaanmetin/go-microservice-starter/internal/infrastructure/database"
	infrahttp "github.com/ozaanmetin/go-microservice-starter/internal/infrastructure/http"
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

	// Setup database connection
	logging.L().Info("Connecting to database...")
	db, err := database.WaitForDB(&cfg.Database, 5, 2*time.Second)
	if err != nil {
		logging.L().WithError(err).Fatal("Failed to connect to database")
	}
	defer database.Close(db)
	logging.L().Info("Database connected successfully")

	// Create HTTP server with route setup from api layer
	server := infrahttp.NewServer(cfg, api.NewRouteSetup(cfg, db))

	// Start server in goroutine
	go func() {
		if err := server.Listen(config.GetServerAddress(cfg)); err != nil {
			logging.L().WithError(err).Fatal("Failed to start server")
		}
	}()

	// Handle graceful shutdown
	gracefulShutdown(server, cfg.Server.ShutdownTimeout)
}

// gracefulShutdown handles graceful shutdown on OS signals
func gracefulShutdown(server *infrahttp.Server, shutdownTimeout time.Duration) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	logging.L().Info("Shutting down server...")

	if err := server.Shutdown(shutdownTimeout); err != nil {
		logging.L().WithError(err).Error("Error during server shutdown")
	}

	logging.L().Info("Server gracefully stopped!")
}
