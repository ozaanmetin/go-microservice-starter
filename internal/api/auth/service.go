package auth

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/ozaanmetin/go-microservice-starter/internal/domain/user"
	pkgJWT "github.com/ozaanmetin/go-microservice-starter/pkg/jwt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotActive      = errors.New("user is not active")
)

// Service handles authentication business logic
type Service struct {
	userRepo   user.Repository
	jwtManager *pkgJWT.Manager
}

// NewService creates a new authentication service
func NewService(userRepo user.Repository, jwtManager *pkgJWT.Manager) *Service {
	return &Service{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// Register creates a new user account
func (s *Service) Register(ctx context.Context, email, password string, firstName, lastName *string) (*user.User, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	newUser := &user.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		FirstName:    firstName,
		LastName:     lastName,
		IsActive:     true,
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

// Login authenticates a user and returns JWT tokens
func (s *Service) Login(ctx context.Context, email, password string) (*pkgJWT.TokenPair, *user.User, error) {
	// Get user by email
	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return nil, nil, ErrInvalidCredentials
		}
		return nil, nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user is active
	if !existingUser.IsActive {
		return nil, nil, ErrUserNotActive
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.PasswordHash), []byte(password)); err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	// Generate JWT tokens
	tokens, err := s.jwtManager.GenerateTokenPair(existingUser.ID, existingUser.Email)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return tokens, existingUser, nil
}

// RefreshToken generates new tokens using a valid refresh token
func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (*pkgJWT.TokenPair, error) {
	// Validate refresh token
	claims, err := s.jwtManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// Verify user still exists and is active
	existingUser, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !existingUser.IsActive {
		return nil, ErrUserNotActive
	}

	// Generate new token pair
	tokens, err := s.jwtManager.GenerateTokenPair(existingUser.ID, existingUser.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return tokens, nil
}

// GetUserByID retrieves a user by ID
func (s *Service) GetUserByID(ctx context.Context, userID int64) (*user.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}
