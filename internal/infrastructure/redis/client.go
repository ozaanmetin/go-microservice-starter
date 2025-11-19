package redis

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ozaanmetin/go-microservice-starter/internal/config"
	"github.com/redis/go-redis/v9"
)

// NewClient creates and returns a new Redis client
func NewClient(cfg *config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + strconv.Itoa(cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Ping to check connection
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return client, nil
}
