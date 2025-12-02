package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/storage/redis/v3"
	appErrors "github.com/ozaanmetin/go-microservice-starter/pkg/errors"
	"github.com/ozaanmetin/go-microservice-starter/pkg/logging"
)

type RateLimiterConfig struct {
	Max          int
	Expiration   time.Duration
	Storage      *redis.Storage
	KeyGenerator KeyGeneratorFunc
	SkipPaths    []string
}

// KeyGeneratorFunc defines the function signature for generating rate limit keys
type KeyGeneratorFunc func(c *fiber.Ctx) string


// Fetches the IP address of the client making the request
func getKeyByIP(c *fiber.Ctx) string {
	return c.IP()
}


// Fetches the user ID from the context for authenticated users
func getKeyByUserId(c *fiber.Ctx) string {
	userID := c.Locals("user_id")

	 // Fallback to IP if user not authenticated
	if userID == nil {
		return c.IP()
	}
	return userID.(string)
}


// Default key generators
var (
	KeyByIP KeyGeneratorFunc = getKeyByIP
	KeyByUserID KeyGeneratorFunc = getKeyByUserId
)


func onLimit(c *fiber.Ctx, keyGen KeyGeneratorFunc) error {
	logging.L().
		WithField("key", keyGen(c)).
		WithField("path", c.Path()).
		Warn("Rate limit exceeded")

	return appErrors.NewTooManyRequestsError("Rate limit exceeded", nil)
}

// RateLimiter middleware limits requests using Redis storage
// Uses IP-based rate limiting by default, but can be customized with KeyGenerator
func RateLimiter(cfg RateLimiterConfig) fiber.Handler {
	// Set default key generator if not provided
	if cfg.KeyGenerator == nil {
		cfg.KeyGenerator = KeyByIP
	}

	return func(c *fiber.Ctx) error {
		// Skip rate limiting for specific paths
		path := c.Path()
		for _, skipPath := range cfg.SkipPaths {
			if path == skipPath {
				return c.Next()
			}
		}
		return limiter.New(limiter.Config{
			Max:          cfg.Max,
			Expiration:   cfg.Expiration,
			KeyGenerator: cfg.KeyGenerator,
			LimitReached: func(c *fiber.Ctx) error {
				return onLimit(c, cfg.KeyGenerator)
			},
			Storage: cfg.Storage,
		})(c)
	}
}

// NewEndpointRateLimiter creates a rate limiter for specific endpoints
// This can be used per-route with custom limits
func NewEndpointRateLimiter(max int, expiration time.Duration, storage *redis.Storage, keyGen KeyGeneratorFunc) fiber.Handler {
	if keyGen == nil {
		keyGen = KeyByIP
	}
	return limiter.New(limiter.Config{
		Max:          max,
		Expiration:   expiration,
		KeyGenerator: keyGen,
		LimitReached: func(c *fiber.Ctx) error {
			return onLimit(c, keyGen)
		},
		Storage: storage,
	})
}
