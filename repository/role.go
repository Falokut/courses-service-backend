package repository

import (
	"context"

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
