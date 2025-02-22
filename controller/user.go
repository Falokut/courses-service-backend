package controller

import (
	"context"
	"courses-service/domain"

	"github.com/Falokut/go-kit/http/apierrors"
	"github.com/Falokut/go-kit/http/types"
	"github.com/pkg/errors"
)

type UserService interface {
	GetRole(ctx context.Context, sessionId string) (*domain.GetRoleResponse, error)
}

type User struct {
	service UserService
}

func NewUser(service UserService) User {
	return User{
		service: service,
	}
}

// GetRole
//
//	@Tags		user
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
//	@Router		/user/get_role [GET]
func (c User) GetRole(ctx context.Context, tokenReq types.BearerToken) (*domain.GetRoleResponse, error) {
	resp, err := c.service.GetRole(ctx, tokenReq.Token)
	switch {
	case errors.Is(err, domain.ErrSessionNotFound):
		return nil, apierrors.NewForbiddenError("forbidden")
	default:
		return resp, err
	}
}
