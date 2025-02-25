package controller

import (
	"context"
	"courses-service/domain"
	"net/http"

	"github.com/Falokut/go-kit/http/apierrors"
	"github.com/Falokut/go-kit/http/types"
	"github.com/pkg/errors"
)

type AuthService interface {
	Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, error)
	Register(ctx context.Context, req domain.RegisterRequest) error
	Logout(ctx context.Context, sessionId string) error
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
		return nil, apierrors.NewForbiddenError(domain.ErrInvalidCredentials.Error())
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

// Logout
//
//	@Tags		user
//	@Summary	Завершить сессию
//	@Produce	json
//
//
//	@Success	200		{object}		any
//	@Failure	401		{object}	apierrors.Error
//	@Failure	500		{object}	apierrors.Error
//
//	@Security	Bearer
//
//	@Router		/auth/logout [POST]
func (c Auth) Logout(ctx context.Context, token types.BearerToken) error {
	return c.service.Logout(ctx, token.Token)
}
