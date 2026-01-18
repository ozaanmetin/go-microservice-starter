package middlewares

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	pkgJWT "github.com/ozaanmetin/go-microservice-starter/pkg/jwt"
)

// UserContextKey is the key used to store user info in context
const UserContextKey = "user"

// AuthMiddleware creates a JWT authentication middleware
// This middleware validates JWT tokens and stores claims in the request context
func AuthMiddleware(jwtManager *pkgJWT.Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    "unauthorized",
				"message": "Missing authorization header",
			})
		}

		// Check Bearer token format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    "unauthorized",
				"message": "Invalid authorization header format",
			})
		}

		tokenString := parts[1]

		// Validate access token
		claims, err := jwtManager.ValidateAccessToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    "unauthorized",
				"message": "Invalid or expired token",
			})
		}

		// Store claims in Fiber context (for framework-level access)
		c.Locals(UserContextKey, claims)

		return c.Next()
	}
}

// GetUserFromContext retrieves user claims from standard context.Context
// This is the framework-agnostic way to access authenticated user information
func GetUserFromContext(ctx context.Context) (*pkgJWT.Claims, bool) {
	claims, ok := ctx.Value(UserContextKey).(*pkgJWT.Claims)
	return claims, ok
}
