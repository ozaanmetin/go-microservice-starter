package middlewares

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	appErrors "github.com/ozaanmetin/go-microservice-starter/pkg/errors"
)

// Recover middleware catches panics and converts them to ServiceError
func Recover() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				var err error
				switch x := r.(type) {
				case string:
					err = fmt.Errorf("%s", x)
				case error:
					err = x
				default:
					err = fmt.Errorf("unknown panic: %v", r)
				}

				// Create ServiceError for panic
				serviceErr := appErrors.NewInternalServerError(err)
				serviceErr.AddDetail("panic_value", r)

				// Initialize empty map if nil to avoid null in JSON
				details := serviceErr.Details
				if details == nil {
					details = make(map[string]interface{})
					serviceErr.Details = details
				}

				// Return ServiceError - will be handled by error handler middleware
				// Logger middleware will log this with all context
				c.Status(fiber.StatusInternalServerError).JSON(serviceErr)
			}
		}()

		return c.Next()
	}
}
