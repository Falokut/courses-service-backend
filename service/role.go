package service

import (
	"context"
	"courses-service/domain"

	"github.com/pkg/errors"
)

type Role struct {
	repo RoleRepo
}

func NewRole(repo RoleRepo) Role {
	return Role{
		repo: repo,
	}
}

func (s Role) GetRoles(ctx context.Context) ([]domain.Role, error) {
	roles, err := s.repo.GetRoles(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "get roles")
	}
	domainRoles := make([]domain.Role, 0, len(roles))
	for _, role := range roles {
		domainRoles = append(domainRoles, domain.Role{
			Id:   role.Id,
			Name: role.Name,
		})
	}
	return domainRoles, nil
}
