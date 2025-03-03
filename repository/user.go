package repository

import (
	"context"
	"database/sql"

	"courses-service/domain"
	"courses-service/entity"

	"github.com/Falokut/go-kit/client/db"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
)

const usernameUniqueConstraint = "users_username_key"

type User struct {
	db db.DB
}

func NewUser(db db.DB) User {
	return User{
		db: db,
	}
}

func (r User) Register(ctx context.Context, user entity.RegisterUser) error {
	const query = `
	INSERT INTO users (username, fio, password, role_id) 
	VALUES($1, $2, $3, $4)
	ON CONFLICT DO NOTHING RETURNING id;`
	var userId int64
	err := r.db.SelectRow(ctx, &userId, query, user.Username, user.Fio, user.PasswordHash, user.RoleId)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return domain.ErrUserAlreadyExists
	case err != nil:
		return errors.WithMessagef(err, "exec query: %s", query)
	}
	return nil
}

func (r User) GetUserByUsername(ctx context.Context, username string) (entity.User, error) {
	const query = `
	SELECT u.id, username, fio, password, role_id, r.name AS role_name
	FROM users u
	JOIN roles r ON u.role_id = r.id 
	WHERE username=$1;`
	var user entity.User
	err := r.db.SelectRow(ctx, &user, query, username)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return entity.User{}, domain.ErrUserNotFound
	case err != nil:
		return entity.User{}, errors.WithMessagef(err, "exec query: %s", query)
	}
	return user, nil
}

func (r User) GetUsers(ctx context.Context, limit int32, offset int32) ([]entity.User, error) {
	const query = `SELECT u.id, username, fio, password, role_id, r.name AS role_name
	FROM users u
	JOIN roles r ON u.role_id = r.id
	WHERE u.id > 0
	ORDER BY u.id
	LIMIT $1 OFFSET $2;`
	var res []entity.User
	err := r.db.Select(ctx, &res, query, limit, offset)
	if err != nil {
		return nil, errors.WithMessagef(err, "exec query: %s", query)
	}
	return res, nil
}

func (r User) GetUserBySessionId(ctx context.Context, sessionId string) (*entity.User, error) {
	const query = `SELECT
		u.id, 
		u.fio,
		u.username,
		u.password, 
		r.id AS role_id,
		r.name AS role_name
	FROM sessions s
	JOIN users u ON s.user_id = u.id
	JOIN roles r ON u.role_id = r.id
	WHERE s.id=$1::TEXT;`
	var user entity.User
	err := r.db.SelectRow(ctx, &user, query, sessionId)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, domain.ErrSessionNotFound
	case err != nil:
		return nil, errors.WithMessagef(err, "exec query: %s", query)
	}
	return &user, nil
}

func (r User) UpsertUser(ctx context.Context, req entity.UpsertUser) error {
	const query = `
	INSERT INTO users (username, fio, password, role_id) VALUES($1, $2, $3, $4)
	ON CONFLICT (username) DO UPDATE
	SET 
		fio = EXCLUDED.fio, 
		password = EXCLUDED.password,
		role_id=EXCLUDED.role_id;`
	_, err := r.db.Exec(ctx, query, req.Username, req.Fio, req.PasswordHash, req.RoleId)
	if err != nil {
		return errors.WithMessagef(err, "exec query: %s", query)
	}
	return nil
}

func (r User) UpdateUser(ctx context.Context, req entity.User) error {
	const query = `
	UPDATE users
	SET username = $1, fio = $2, password = $3, role_id = $4
	WHERE id = $5;`
	_, err := r.db.Exec(ctx, query, req.Username, req.Fio, req.Password, req.RoleId, req.Id)
	pgError := &pgconn.PgError{}
	switch {
	case errors.As(err, &pgError) && pgError.ConstraintName == usernameUniqueConstraint:
		return domain.ErrUserAlreadyExists
	case err != nil:
		return errors.WithMessagef(err, "exec query: %s", query)
	default:
		return nil
	}
}

func (r User) DeleteUser(ctx context.Context, userId int32) error {
	const query = "DELETE FROM users WHERE id=$1"
	_, err := r.db.Exec(ctx, query, userId)
	if err != nil {
		return errors.WithMessagef(err, "exec query: %s", query)
	}
	return nil
}

func (r User) GetUsersByRoleName(ctx context.Context, roleName string) ([]entity.User, error) {
	const query = `SELECT u.id, username, fio, password, role_id, r.name AS role_name
	FROM users u
	JOIN roles r ON u.role_id = r.id
	WHERE r.name =$1
	ORDER BY u.id;`
	var res []entity.User
	err := r.db.Select(ctx, &res, query, roleName)
	if err != nil {
		return nil, errors.WithMessagef(err, "exec query: %s", query)
	}
	return res, nil
}
