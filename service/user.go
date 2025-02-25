package service

import (
	"context"
	"courses-service/domain"

	"github.com/pkg/errors"
)

type User struct {
	authRepo AuthRepo
	userRepo UserRepo
}

func NewUser(authRepo AuthRepo, userRepo UserRepo) User {
	return User{
		authRepo: authRepo,
		userRepo: userRepo,
	}
}

func (s User) GetRole(ctx context.Context, sessionId string) (*domain.GetRoleResponse, error) {
	userSession, err := s.authRepo.GetUserSession(ctx, sessionId)
	if err != nil {
		return nil, errors.WithMessage(err, "get user session")
	}
	return &domain.GetRoleResponse{
		RoleName: userSession.RoleName,
	}, nil
}

func (s User) GetUsers(ctx context.Context, req domain.LimitOffsetRequest) ([]domain.User, error) {
	users, err := s.userRepo.GetUsers(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, errors.WithMessage(err, "get roles")
	}
	domainUsers := make([]domain.User, 0, len(users))
	for _, user := range users {
		domainUsers = append(domainUsers,
			domain.User{
				Id:       user.Id,
				Username: user.Username,
				Fio:      user.Fio,
				RoleName: user.RoleName,
				RoleId:   user.RoleId,
			})
	}
	return domainUsers, nil
}

func (s User) DeleteUser(ctx context.Context, userId int32) error {
	err := s.userRepo.DeleteUser(ctx, userId)
	if err != nil {
		return errors.WithMessage(err, "get user session")
	}
	return nil
}
