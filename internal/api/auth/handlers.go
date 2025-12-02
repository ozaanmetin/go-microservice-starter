package auth

import (
	"context"
	"errors"

	"github.com/ozaanmetin/go-microservice-starter/internal/domain/user"
	appErrors "github.com/ozaanmetin/go-microservice-starter/pkg/errors"
	pkgJWT "github.com/ozaanmetin/go-microservice-starter/pkg/jwt"
)

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	Email     string  `json:"email" validate:"required,email"`
	Password  string  `json:"password" validate:"required,min=8"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
}

// RegisterResponse represents the registration response
type RegisterResponse struct {
	User   *UserResponse     `json:"user"`
	Tokens *pkgJWT.TokenPair `json:"tokens"`
}

func (r *RegisterResponse) StatusCode() int {
	return 201
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	User   *UserResponse     `json:"user"`
	Tokens *pkgJWT.TokenPair `json:"tokens"`
}

// RefreshTokenRequest represents the refresh token request payload
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// RefreshTokenResponse represents the refresh token response
type RefreshTokenResponse struct {
	Tokens *pkgJWT.TokenPair `json:"tokens"`
}

// UserResponse represents user data in responses
type UserResponse struct {
	ID        int64   `json:"id"`
	Email     string  `json:"email"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	IsActive  bool    `json:"is_active"`
}

// toUserResponse converts a user entity to response format
func toUserResponse(u *user.User) *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		IsActive:  u.IsActive,
	}
}

// RegisterHandler handles user registration
type RegisterHandler struct {
	service *Service
}

// NewRegisterHandler creates a new register handler
func NewRegisterHandler(service *Service) *RegisterHandler {
	return &RegisterHandler{service: service}
}

// Handle processes the registration request
func (h *RegisterHandler) Handle(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	// Register user
	newUser, err := h.service.Register(ctx, req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		if errors.Is(err, user.ErrUserAlreadyExists) {
			return nil, appErrors.NewConflictError("User with this email already exists", err)
		}
		return nil, appErrors.NewInternalServerError(err)
	}

	// Generate tokens
	tokens, err := h.service.jwtManager.GenerateTokenPair(newUser.ID, newUser.Email)
	if err != nil {
		return nil, appErrors.NewInternalServerError(err)
	}

	return &RegisterResponse{
		User:   toUserResponse(newUser),
		Tokens: tokens,
	}, nil
}

// LoginHandler handles user login
type LoginHandler struct {
	service *Service
}

// NewLoginHandler creates a new login handler
func NewLoginHandler(service *Service) *LoginHandler {
	return &LoginHandler{service: service}
}

// Handle processes the login request
func (h *LoginHandler) Handle(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	tokens, existingUser, err := h.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			return nil, appErrors.NewUnauthorizedError("Invalid email or password", err)
		}
		if errors.Is(err, ErrUserNotActive) {
			return nil, appErrors.NewForbiddenError("User account is not active", err)
		}
		return nil, appErrors.NewInternalServerError(err)
	}

	return &LoginResponse{
		User:   toUserResponse(existingUser),
		Tokens: tokens,
	}, nil
}

// RefreshTokenHandler handles token refresh
type RefreshTokenHandler struct {
	service *Service
}

// NewRefreshTokenHandler creates a new refresh token handler
func NewRefreshTokenHandler(service *Service) *RefreshTokenHandler {
	return &RefreshTokenHandler{service: service}
}

// Handle processes the token refresh request
func (h *RefreshTokenHandler) Handle(ctx context.Context, req *RefreshTokenRequest) (*RefreshTokenResponse, error) {
	tokens, err := h.service.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		if errors.Is(err, pkgJWT.ErrExpiredToken) {
			return nil, appErrors.NewUnauthorizedError("Refresh token has expired", err)
		}
		if errors.Is(err, pkgJWT.ErrInvalidToken) || errors.Is(err, pkgJWT.ErrInvalidSignature) {
			return nil, appErrors.NewUnauthorizedError("Invalid refresh token", err)
		}
		if errors.Is(err, ErrUserNotActive) {
			return nil, appErrors.NewForbiddenError("User account is not active", err)
		}
		return nil, appErrors.NewInternalServerError(err)
	}

	return &RefreshTokenResponse{
		Tokens: tokens,
	}, nil
}
