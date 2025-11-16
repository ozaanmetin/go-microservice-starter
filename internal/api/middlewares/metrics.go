package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ozaanmetin/go-microservice-starter/pkg/metrics"
)

// Metrics middleware records HTTP request metrics
func Metrics() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip metrics endpoint to avoid recursion
		if c.Path() == "/metrics" {
			return c.Next()
		}

		// Increment in-flight requests
		metrics.IncrementInFlight()
		defer metrics.DecrementInFlight()

		// Record start time
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Get status code
		statusCode := c.Response().StatusCode()

		// Record metrics
		metrics.RecordHTTPRequest(
			c.Method(),
			c.Path(),
			statusCode,
			duration,
		)

		return err
	}
}
