package controller

import (
	"context"
	"courses-service/domain"
	_ "github.com/Falokut/go-kit/http/apierrors"
)

type RoleService interface {
	GetRoles(ctx context.Context) ([]domain.Role, error)
}

type Role struct {
	service RoleService
}

func NewRole(service RoleService) Role {
	return Role{
		service: service,
	}
}

// GetRoles
//
//	@Tags		role
//	@Summary	Получить список ролей
//	@Accept		json
//	@Produce	json
//
//
//	@Success	200	{array}		domain.Role
//	@Failure	404	{object}	apierrors.Error
//	@Failure	500	{object}	apierrors.Error
//	@Router		/roles [GET]
func (c Role) GetRoles(ctx context.Context) ([]domain.Role, error) {
	return c.service.GetRoles(ctx)
}
