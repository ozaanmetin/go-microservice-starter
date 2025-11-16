package middlewares

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	appErrors "github.com/ozaanmetin/go-microservice-starter/pkg/errors"
)

type ServiceErrorResponse struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details"`
}

// ServiceErrorErrorHandler returns a Fiber error handler that handles ServiceError
func ServiceErrorErrorHandler() fiber.ErrorHandler {
	// Get Fiber's default error handler
	defaultHandler := fiber.DefaultErrorHandler

	return func(c *fiber.Ctx, err error) error {
		// Check if it's a ServiceError
		var serviceErr *appErrors.ServiceError

		if !errors.As(err, &serviceErr) {
			// Not a ServiceError - use Fiber's default handler
			return defaultHandler(c, err)
		}

		// It's a ServiceError - handle it
		statusCode := serviceErr.StatusCode

		// Build response
		response := ServiceErrorResponse{
			Code:    serviceErr.Code,
			Message: serviceErr.Message,
			Details: serviceErr.Details,
		}

		// Initialize empty map if nil to avoid null in JSON
		if response.Details == nil {
			response.Details = make(map[string]interface{})
		}
		return c.Status(statusCode).JSON(response)
	}
}
