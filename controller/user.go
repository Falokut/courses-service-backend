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
	GetUsers(ctx context.Context, req domain.LimitOffsetRequest) ([]domain.User, error)
	DeleteUser(ctx context.Context, userId int32) error
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
//	@Produce	json
//
//
//	@Success	200	{object}	domain.GetRoleResponse
//	@Failure	401	{object}	apierrors.Error
//	@Failure	403	{object}	apierrors.Error
//	@Failure	500	{object}	apierrors.Error
//
//	@Security	Bearer
//
//	@Router		/users/get_role [GET]
func (c User) GetRole(ctx context.Context, tokenReq types.BearerToken) (*domain.GetRoleResponse, error) {
	resp, err := c.service.GetRole(ctx, tokenReq.Token)
	switch {
	case errors.Is(err, domain.ErrSessionNotFound):
		return nil, apierrors.NewForbiddenError(domain.ErrForbidden.Error())
	default:
		return resp, err
	}
}

// GetUsers
//
//	@Tags		user
//	@Summary	Получить пользователей
//	@Produce	json
//
//	@Param		limit	query		int	true	"максимальное количество записей"
//	@Param		offset	query		int	true	"отступ поиска записей"
//
//	@Success	200		{array}		domain.User
//	@Failure	400		{object}	apierrors.Error
//	@Failure	401		{object}	apierrors.Error
//	@Failure	403		{object}	apierrors.Error
//	@Failure	500		{object}	apierrors.Error
//
//	@Security	Bearer
//
//	@Router		/users [GET]
func (c User) GetUsers(ctx context.Context, req domain.LimitOffsetRequest) ([]domain.User, error) {
	return c.service.GetUsers(ctx, req)
}

// DeleteUser
//
//	@Tags		user
//	@Summary	Удалить пользователя
//	@Produce	json
//
//	@Param		userId	path		int	true	"максимальное количество записей"
//
//	@Success	200		{object}		any
//	@Failure	400		{object}	apierrors.Error
//	@Failure	401		{object}	apierrors.Error
//	@Failure	403		{object}	apierrors.Error
//	@Failure	500		{object}	apierrors.Error
//
//	@Security	Bearer
//
//	@Router		/users/{userId} [DELETE]
func (c User) DeleteUser(ctx context.Context, req domain.DeleteUserRequest) error {
	return c.service.DeleteUser(ctx, req.UserId)
}
