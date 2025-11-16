package http

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	appErrors "github.com/ozaanmetin/go-microservice-starter/pkg/errors"
)

// Request is a marker interface for all request types
type Request interface{}

// Response is a marker interface for all response types
type Response interface{}

// StatusCodeProvider allows responses to specify their HTTP status code
type StatusCodeProvider interface {
	StatusCode() int
}

// HandlerInterface defines the generic handler interface for business logic
type HandlerInterface[R Request, Res Response] interface {
	Handle(ctx context.Context, req *R) (*Res, error)
}

// AdaptHandler adapts a generic HandlerInterface to Fiber's handler signature
// It handles request parsing (body, params, query, headers) and error responses
func AdaptHandler[R Request, Res Response](handler HandlerInterface[R, Res]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req R

		// Parse body
		if err := c.BodyParser(&req); err != nil && err != fiber.ErrUnprocessableEntity {
			return appErrors.NewBadRequestError("Invalid request body", err)
		}

		// Parse params
		if err := c.ParamsParser(&req); err != nil {
			return appErrors.NewBadRequestError("Invalid URL parameters", err)
		}

		// Parse query
		if err := c.QueryParser(&req); err != nil {
			return appErrors.NewBadRequestError("Invalid query parameters", err)
		}

		// Parse headers
		if err := c.ReqHeaderParser(&req); err != nil {
			return appErrors.NewBadRequestError("Invalid headers", err)
		}

		// Get user context
		ctx := c.UserContext()

		// Call business handler
		res, err := handler.Handle(ctx, &req)

		if err != nil {
			// Check if it's already a ServiceError
			var serviceErr *appErrors.ServiceError
			if errors.As(err, &serviceErr) {
				return err // Return as-is, error handler middleware will handle it
			}

			// Wrap unknown errors as internal server error
			return appErrors.NewInternalServerError(err)
		}

		// Check if response implements StatusCodeProvider
		if provider, ok := any(res).(StatusCodeProvider); ok {
			return c.Status(provider.StatusCode()).JSON(res)
		}

		// Default to 200 OK
		return c.JSON(res)
	}
}
