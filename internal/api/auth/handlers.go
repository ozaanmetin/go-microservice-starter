package auth

import (
	"context"
	"errors"

	"github.com/ozaanmetin/go-microservice-starter/internal/domain/user"
	appErrors "github.com/ozaanmetin/go-microservice-starter/pkg/errors"
	pkgJWT "github.com/ozaanmetin/go-microservice-starter/pkg/jwt"
)

// Register related structs

type RegisterRequest struct {
	Email     string  `json:"email" validate:"required,email"`
	Password  string  `json:"password" validate:"required,min=8"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
}

type RegisterResponse struct {
	User   *UserResponse     `json:"user"`
	Tokens *pkgJWT.TokenPair `json:"tokens"`
}

func (r *RegisterResponse) StatusCode() int {
	return 201
}


// Login related structs

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	User   *UserResponse     `json:"user"`
	Tokens *pkgJWT.TokenPair `json:"tokens"`
}

// Refresh related structs

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshTokenResponse struct {
	Tokens *pkgJWT.TokenPair `json:"tokens"`
}


// User related structs

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


// Register Handler handles user registration

type RegisterHandler struct {
	service *AuthService
}

func NewRegisterHandler(service *AuthService) *RegisterHandler {
	return &RegisterHandler{service: service}
}

func (h *RegisterHandler) Handle(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	// register the user
	newUser, err := h.service.Register(ctx, req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		if errors.Is(err, user.ErrUserAlreadyExists) {
			return nil, appErrors.NewConflictError("User with this email already exists", err)
		}
		return nil, appErrors.NewInternalServerError(err)
	}

	// generates tokens
	tokens, err := h.service.jwtManager.GenerateTokenPair(newUser.ID, newUser.Email)
	if err != nil {
		return nil, appErrors.NewInternalServerError(err)
	}

	return &RegisterResponse{
		User:   toUserResponse(newUser),
		Tokens: tokens,
	}, nil
}


// Login Handler handles user login

type LoginHandler struct {
	service *AuthService
}

func NewLoginHandler(service *AuthService) *LoginHandler {
	return &LoginHandler{service: service}
}

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


// Refresh Token Handler handles token refreshing

type RefreshTokenHandler struct {
	service *AuthService
}

func NewRefreshTokenHandler(service *AuthService) *RefreshTokenHandler {
	return &RefreshTokenHandler{service: service}
}

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
