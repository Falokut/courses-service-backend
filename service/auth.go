package service

import (
	"context"
	"courses-service/conf"
	"courses-service/domain"
	"courses-service/entity"

	"github.com/Falokut/go-kit/http/apierrors"
	"github.com/Falokut/go-kit/jwt"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepo interface {
	GetUserByUsername(ctx context.Context, username string) (entity.User, error)
	Register(ctx context.Context, req entity.RegisterUser) error
}

type Auth struct {
	cfg  conf.Auth
	repo AuthRepo
}

func NewAuth(cfg conf.Auth, repo AuthRepo) Auth {
	return Auth{
		cfg:  cfg,
		repo: repo,
	}
}

func (s Auth) Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, error) {
	user, err := s.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.WithMessage(err, "get user by username")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	tokenValue := entity.TokenUserInfo{
		UserId:   user.Id,
		RoleName: user.RoleName,
	}
	accessToken, err := jwt.GenerateToken(
		s.cfg.Access.Secret,
		s.cfg.Access.TTL,
		&tokenValue,
	)
	if err != nil {
		return nil, errors.WithMessage(err, "generate access token")
	}

	refreshToken, err := jwt.GenerateToken(
		s.cfg.Refresh.Secret,
		s.cfg.Refresh.TTL,
		&tokenValue,
	)
	if err != nil {
		return nil, errors.WithMessage(err, "generate refresh token")
	}

	return &domain.LoginResponse{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}, nil
}

func (s Auth) Register(ctx context.Context, req domain.RegisterRequest) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), s.cfg.BcryptCost)
	if err != nil {
		return errors.WithMessage(err, "generate from passport")
	}

	err = s.repo.Register(ctx, entity.RegisterUser{
		Username:     req.Username,
		Fio:          req.Fio,
		PasswordHash: string(passwordHash),
		RoleId:       req.RoleId,
	})
	if err != nil {
		return errors.WithMessage(err, "register")
	}
	return nil
}

func (s Auth) RefreshAccessToken(ctx context.Context, refreshToken string) (*jwt.TokenResponse, error) {
	tokenValue := entity.TokenUserInfo{}
	err := jwt.ParseToken(refreshToken, s.cfg.Refresh.Secret, &tokenValue)
	if err != nil {
		return nil, errors.WithMessage(err, "parse token")
	}
	accessToken, err := jwt.GenerateToken(
		s.cfg.Access.Secret,
		s.cfg.Access.TTL,
		&tokenValue,
	)
	if err != nil {
		return nil, errors.WithMessage(err, "generate access token")
	}
	return accessToken, nil
}

func (s Auth) GetRole(ctx context.Context, accessToken string) (*domain.GetRoleResponse, error) {
	tokenValue := entity.TokenUserInfo{}
	err := jwt.ParseToken(accessToken, s.cfg.Access.Secret, &tokenValue)
	if err != nil {
		return nil, apierrors.NewUnauthorizedError("invalid token") //nolint:wrapcheck
	}

	return &domain.GetRoleResponse{
		RoleName: tokenValue.RoleName,
	}, nil
}
