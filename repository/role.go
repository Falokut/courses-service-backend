package repository

import (
	"context"
	"courses-service/entity"

	"github.com/Falokut/go-kit/client/db"
	"github.com/pkg/errors"
)

type Role struct {
	db db.DB
}

func NewRole(db db.DB) Role {
	return Role{
		db: db,
	}
}

func (r Role) GetRoleId(ctx context.Context, roleName string) (int32, error) {
	const query = "SELECT id FROM roles WHERE name = $1;"
	var roleId int32
	err := r.db.SelectRow(ctx, &roleId, query, roleName)
	if err != nil {
		return -1, errors.WithMessagef(err, "exec query: %s", query)
	}
	return roleId, nil
}

func (r Role) GetRoles(ctx context.Context) ([]entity.Role, error) {
	const query = "SELECT id, name FROM roles;"
	roles := make([]entity.Role, 0)
	err := r.db.Select(ctx, &roles, query)
	if err != nil {
		return nil, errors.WithMessagef(err, "exec query: %s", query)
	}
	return roles, nil
}
