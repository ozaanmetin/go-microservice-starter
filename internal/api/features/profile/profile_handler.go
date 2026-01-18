package profile

import (
	"context"

	"github.com/ozaanmetin/go-microservice-starter/internal/middlewares"
	appErrors "github.com/ozaanmetin/go-microservice-starter/pkg/errors"
	pkgJWT "github.com/ozaanmetin/go-microservice-starter/pkg/jwt"
)

// GetProfileRequest represents the profile request
// This handler gets the user info from the JWT token via middleware
type GetProfileRequest struct{}

// GetProfileResponse represents the profile response
type GetProfileResponse struct {
	UserID    int64  `json:"user_id"`
	Email     string `json:"email"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
}

// GetProfileHandler handles fetching user profile
type GetProfileHandler struct {
	profileService *ProfileService
}

// NewGetProfileHandler creates a new profile handler
func NewGetProfileHandler(profileService *ProfileService) *GetProfileHandler {
	return &GetProfileHandler{
		profileService: profileService,
	}
}

// Handle processes the get profile request
// This demonstrates how to use JWT authentication in a protected endpoint
// The handler is completely framework-agnostic and uses standard context.Context
func (h *GetProfileHandler) Handle(ctx context.Context, req *GetProfileRequest) (*GetProfileResponse, error) {
	// Get user claims from context (set by AuthMiddleware and transferred by adapter)
	claims, ok := middlewares.GetUserFromContext(ctx)
	if !ok {
		return nil, appErrors.NewUnauthorizedError("User not authenticated", nil)
	}

	// Optionally, you can fetch full user details from database
	user, err := h.profileService.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, appErrors.NewInternalServerError(err)
	}

	// Return profile information
	return &GetProfileResponse{
		UserID: user.ID,
		Email:  user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, nil
}

// ExtractUserClaims is a helper function to extract user claims from context
// This can be used in any handler that needs to access the authenticated user
// It's completely framework-agnostic
func ExtractUserClaims(ctx context.Context) (*pkgJWT.Claims, error) {
	claims, ok := middlewares.GetUserFromContext(ctx)
	if !ok {
		return nil, appErrors.NewUnauthorizedError("User not authenticated", nil)
	}

	return claims, nil
}
