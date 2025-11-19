package redis

import (
	redisstorage "github.com/gofiber/storage/redis/v3"
	"github.com/ozaanmetin/go-microservice-starter/internal/config"
)

// NewStorage creates a new Redis storage for Fiber middleware (rate limiting, session, etc.)
func NewStorage(cfg *config.RedisConfig) *redisstorage.Storage {
	return redisstorage.New(redisstorage.Config{
		Host:     cfg.Host,
		Port:     cfg.Port,
		Password: cfg.Password,
		Database: cfg.DB,
	})
}
