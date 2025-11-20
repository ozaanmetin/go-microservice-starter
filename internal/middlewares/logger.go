package middlewares

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ozaanmetin/go-microservice-starter/pkg/logging"
	appErrors "github.com/ozaanmetin/go-microservice-starter/pkg/errors"
)

type LoggerConfig struct {
	// Skip logging for specific paths (e.g., /metrics, /health for noise reduction)
	SkipPaths []string
}

// Logger middleware logs all HTTP requests with details
func Logger(cfg LoggerConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Path()
		for _, skipPath := range cfg.SkipPaths {
			if path == skipPath {
				return c.Next()
			}
		}

		start := time.Now()
		requestID := GetRequestID(c)

		// Process request and capture error
		err := c.Next()

		// Log after the chain completes
		duration := time.Since(start)

		// Determine the status code:
		// If there's an error, extract status from ServiceError or Fiber error before error handler runs
		// Otherwise, get it from response (for successful requests)
		statusCode := c.Response().StatusCode()

		if err != nil {
			// Check if it's a ServiceError and extract status code
			statusCode = getStatusCodeFromError(err)
		}

		logger := logging.L().
			WithField("method", c.Method()).
			WithField("path", path).
			WithField("status", statusCode).
			WithField("duration_ms", duration.Milliseconds()).
			WithField("ip", c.IP()).
			WithField("user_agent", c.Get("User-Agent")).
			WithField("bytes_sent", len(c.Response().Body())).
			WithField("is_successful", statusCode >= 200 && statusCode < 400)

		if requestID != "" {
			logger = logger.WithField("request_id", requestID)
		}

		if err != nil {
			logger = logger.WithError(err)
		}

		logForRequest(logger, statusCode)
		return err
	}
}

// Extract status code from error if it's a ServiceError or Fiber error
func getStatusCodeFromError(err error) int {
	var serviceErr *appErrors.ServiceError
	if errors.As(err, &serviceErr) {
		return serviceErr.StatusCode
	} else if fiberErr, ok := err.(*fiber.Error); ok {
		return fiberErr.Code
	}
	return fiber.StatusInternalServerError
}

// Log at appropriate level based on status code
func logForRequest(l *logging.Logger, statusCode int) {
	switch {
	case statusCode >= 500:
		l.Error("HTTP request completed with server error")
	case statusCode >= 400:
		l.Warn("HTTP request completed with client error")
	default:
		l.Info("HTTP request completed")
	}
}
