package profile

import (
	"context"	
	"github.com/ozaanmetin/go-microservice-starter/internal/domain/user"
)

// ProfileService provides profile related operations

type ProfileService struct {
	userRepo user.Repository
}

func NewProfileService(userRepo user.Repository) *ProfileService{
	return &ProfileService{
		userRepo: userRepo,
	}
}

func (s *ProfileService) GetUserByID(ctx context.Context, userID int64) (*user.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}