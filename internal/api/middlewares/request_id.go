package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

// GetRequestID retrieves the request ID from fiber context
// Uses Fiber's built-in requestid middleware context key
func GetRequestID(c *fiber.Ctx) string {
	// Fiber's requestid middleware stores ID in locals with "requestid" key
	if id, ok := c.Locals("requestid").(string); ok {
		return id
	}
	return ""
}

// RequestID returns Fiber's built-in request ID middleware
// It generates a unique ID for each request and adds it to X-Request-Id header
func RequestID() fiber.Handler {
	return requestid.New(requestid.Config{
		Header: "X-Request-Id",
		// Generator can be customized here if needed
		// Generator: func() string { return uuid.New().String() },
	})
}
