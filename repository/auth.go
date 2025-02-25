package repository

import (
	"context"
	"database/sql"
	"time"

	"courses-service/domain"
	"courses-service/entity"

	"github.com/Falokut/go-kit/client/db"
	"github.com/pkg/errors"
)

type Auth struct {
	db db.DB
}

func NewAuth(db db.DB) Auth {
	return Auth{
		db: db,
	}
}

func (r Auth) GetUserSession(ctx context.Context, sessionId string) (entity.UserSession, error) {
	const query = `SELECT s.id, s.user_id, s.created_at, r.name AS role_name
	FROM sessions s
	JOIN users u ON s.user_id = u.id
	JOIN roles r ON u.role_id = r.id
	WHERE s.id=$1::TEXT;`
	var userSession entity.UserSession
	err := r.db.SelectRow(ctx, &userSession, query, sessionId)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return entity.UserSession{}, domain.ErrSessionNotFound
	case err != nil:
		return entity.UserSession{}, errors.WithMessagef(err, "exec query: %s", query)
	}
	return userSession, nil
}

func (r Auth) InsertSession(ctx context.Context, session entity.Session) error {
	const query = "INSERT INTO sessions (id, user_id, created_at) VALUES($1, $2, $3);"
	_, err := r.db.Exec(ctx, query, session.Id, session.UserId, session.CreatedAt)
	if err != nil {
		return errors.WithMessagef(err, "exec query: %s", query)
	}
	return nil
}

func (r Auth) DeleteExpiredSessions(ctx context.Context, startFrom time.Time) error {
	const query = "DELETE FROM sessions WHERE created_at <= $1;"
	_, err := r.db.Exec(ctx, query, startFrom)
	if err != nil {
		return errors.WithMessagef(err, "exec query: %s", query)
	}
	return nil
}

func (r Auth) DeleteSession(ctx context.Context, sessionId string) error {
	const query = "DELETE FROM sessions WHERE id = $1;"
	_, err := r.db.Exec(ctx, query, sessionId)
	if err != nil {
		return errors.WithMessagef(err, "exec query: %s", query)
	}
	return nil
}
