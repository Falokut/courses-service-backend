package controller

import (
	"context"
	"courses-service/domain"
	"net/http"

	"github.com/Falokut/go-kit/http/apierrors"
	"github.com/Falokut/go-kit/http/types"
	"github.com/Falokut/go-kit/jwt"
	"github.com/pkg/errors"
)

type AuthService interface {
	Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, error)
	Register(ctx context.Context, req domain.RegisterRequest) error
	RefreshAccessToken(ctx context.Context, refreshToken string) (*jwt.TokenResponse, error)
	GetRole(ctx context.Context, accessToken string) (*domain.GetRoleResponse, error)
}

type Auth struct {
	service AuthService
}

func NewAuth(service AuthService) Auth {
	return Auth{
		service: service,
	}
}

// Login
//
//	@Tags		auth
//	@Summary	Войти в аккаунт
//	@Accept		json
//	@Produce	json
//
//	@Param		body	body		domain.LoginRequest	true	"тело запроса"
//
//	@Success	200		{object}	domain.LoginResponse
//	@Failure	404		{object}	apierrors.Error
//	@Failure	500		{object}	apierrors.Error
//	@Router		/auth/login [POST]
func (c Auth) Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, error) {
	tokens, err := c.service.Login(ctx, req)
	switch {
	case errors.Is(err, domain.ErrUserNotFound), errors.Is(err, domain.ErrInvalidCredentials):
		return nil, apierrors.NewForbiddenError("invalid credentials")
	default:
		return tokens, err
	}
}

// Register
//
//	@Tags		auth
//	@Summary	Зарегистрировать пользователя
//	@Accept		json
//	@Produce	json
//
//	@Param		body	body		domain.RegisterRequest	true	"тело запроса"
//
//	@Success	200		{object}	any
//	@Failure	404		{object}	apierrors.Error
//	@Failure	500		{object}	apierrors.Error
//	@Router		/auth/register [POST]
func (c Auth) Register(ctx context.Context, req domain.RegisterRequest) error {
	err := c.service.Register(ctx, req)
	switch {
	case errors.Is(err, domain.ErrUserAlreadyExists):
		return apierrors.New(
			http.StatusConflict,
			domain.ErrCodeUserAlreadyExists,
			domain.ErrUserAlreadyExists.Error(),
			err,
		)
	default:
		return err
	}
}

// RefreshAccessToken
//
//	@Tags		auth
//	@Summary	Обновить токен доступа
//	@Accept		json
//	@Produce	json
//
//
//	@Success	200	{object}	jwt.TokenResponse
//	@Failure	404	{object}	apierrors.Error
//	@Failure	401	{object}	apierrors.Error
//	@Failure	500	{object}	apierrors.Error
//
//	@Security	Bearer
//
//	@Router		/auth/access_token [GET]
func (c Auth) RefreshAccessToken(ctx context.Context, tokenReq types.BearerToken) (*jwt.TokenResponse, error) {
	return c.service.RefreshAccessToken(ctx, tokenReq.Token)
}

// GetRole
//
//	@Tags		auth
//	@Summary	Получить роль пользователя
//	@Accept		json
//	@Produce	json
//
//
//	@Success	200	{object}	domain.GetRoleResponse
//	@Failure	404	{object}	apierrors.Error
//	@Failure	401	{object}	apierrors.Error
//	@Failure	500	{object}	apierrors.Error
//
//	@Security	Bearer
//
//	@Router		/auth/get_user_role [GET]
func (c Auth) GetRole(ctx context.Context, tokenReq types.BearerToken) (*domain.GetRoleResponse, error) {
	return c.service.GetRole(ctx, tokenReq.Token)
}
