package service

import (
	"context"
	"courses-service/domain"
	"courses-service/entity"

	"github.com/pkg/errors"
)

type AuthRepo interface {
	GetUserSession(ctx context.Context, sessionId string) (entity.UserSession, error)
}

type User struct {
	userRepo AuthRepo
}

func NewUser(userRepo AuthRepo) User {
	return User{
		userRepo: userRepo,
	}
}

func (s User) GetRole(ctx context.Context, sessionId string) (*domain.GetRoleResponse, error) {
	userSession, err := s.userRepo.GetUserSession(ctx, sessionId)
	if err != nil {
		return nil, errors.WithMessage(err, "get user session")
	}
	return &domain.GetRoleResponse{
		RoleName: userSession.RoleName,
	}, nil
}
