package repository

import (
	"context"
	"database/sql"

	"courses-service/domain"
	"courses-service/entity"

	"github.com/Falokut/go-kit/client/db"
	"github.com/pkg/errors"
)

type User struct {
	cli *db.Client
}

func NewUser(cli *db.Client) User {
	return User{
		cli: cli,
	}
}

func (r User) Register(ctx context.Context, user entity.RegisterUser) error {
	query := `
	INSERT INTO users (username, fio, password, role_id) 
	VALUES($1, $2, $3, $4)
	ON CONFLICT DO NOTHING RETURNING id;`
	var userId int64
	err := r.cli.GetContext(ctx, &userId, query, user.Username, user.Fio, user.PasswordHash, user.RoleId)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return domain.ErrUserAlreadyExists
	case err != nil:
		return errors.WithMessagef(err, "exec query: %s", query)
	}
	return nil
}

func (r User) GetUserByUsername(ctx context.Context, username string) (entity.User, error) {
	query := `
	SELECT u.id, username, fio, password, role_id, r.name AS role_name
	FROM users u
	JOIN roles r ON u.role_id = r.id 
	WHERE username=$1;`
	var user entity.User
	err := r.cli.GetContext(ctx, &user, query, username)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return entity.User{}, domain.ErrUserNotFound
	case err != nil:
		return entity.User{}, errors.WithMessagef(err, "exec query: %s", query)
	}
	return user, nil
}

func (r User) GetUsers(ctx context.Context, limit int32, offset int32) ([]entity.User, error) {
	query := "SELECT id, username, fio, role_id FROM users LIMIT $1 OFFSET $2"
	var res []entity.User
	err := r.cli.SelectContext(ctx, &res, query, limit, offset)
	if err != nil {
		return nil, errors.WithMessagef(err, "exec query: %s", query)
	}
	return res, nil
}
