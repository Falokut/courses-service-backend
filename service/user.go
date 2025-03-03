package service

import (
	"context"
	"courses-service/conf"
	"courses-service/domain"
	"courses-service/entity"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	cfg      conf.Auth
	authRepo AuthRepo
	userRepo UserRepo
}

func NewUser(cfg conf.Auth, authRepo AuthRepo, userRepo UserRepo) User {
	return User{
		cfg:      cfg,
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
		return nil, errors.WithMessage(err, "get users")
	}
	return entityUserToDomain(users), nil
}

func (s User) GetUserProfile(ctx context.Context, sessionId string) (*domain.UserProfile, error) {
	user, err := s.userRepo.GetUserBySessionId(ctx, sessionId)
	if err != nil {
		return nil, errors.WithMessage(err, "get user by session id")
	}
	return &domain.UserProfile{
		Id:       user.Id,
		Username: user.Username,
		Fio:      user.Fio,
		RoleName: user.RoleName,
		RoleId:   user.RoleId,
	}, nil
}

func (s User) DeleteUser(ctx context.Context, userId int32) error {
	err := s.userRepo.DeleteUser(ctx, userId)
	if err != nil {
		return errors.WithMessage(err, "get user session")
	}
	return nil
}

func (s User) EditUser(ctx context.Context, req domain.EditUserRequest) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), s.cfg.BcryptCost)
	if err != nil {
		return errors.WithMessage(err, "generate from passport")
	}

	err = s.userRepo.UpdateUser(ctx, entity.User{
		Id:       req.UserId,
		Username: req.Username,
		Fio:      req.Fio,
		Password: string(passwordHash),
		RoleId:   req.RoleId,
	})
	if err != nil {
		return errors.WithMessage(err, "update user")
	}
	return nil
}

func (s User) GetLectors(ctx context.Context) ([]domain.User, error) {
	users, err := s.userRepo.GetUsersByRoleName(ctx, domain.LectorType)
	if err != nil {
		return nil, errors.WithMessage(err, "get users by role name")
	}
	return entityUserToDomain(users), nil
}

func entityUserToDomain(users []entity.User) []domain.User {
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
	return domainUsers
}
