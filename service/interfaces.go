package service

import (
	"context"
	"courses-service/entity"
)

type UserRepo interface {
	Register(ctx context.Context, req entity.RegisterUser) error
	UpsertUser(ctx context.Context, req entity.UpsertUser) error
	GetUsers(ctx context.Context, limit int32, offset int32) ([]entity.User, error)
	DeleteUser(ctx context.Context, userId int32) error
}

type AuthTxRunner interface {
	LoginTransaction(ctx context.Context, txFunc func(ctx context.Context, tx LoginTx) error) error
}

type LoginTx interface {
	GetUserByUsername(ctx context.Context, username string) (entity.User, error)
	InsertSession(ctx context.Context, session entity.Session) error
}

type AuthRepo interface {
	GetUserSession(ctx context.Context, sessionId string) (entity.UserSession, error)
	DeleteSession(ctx context.Context, sessionId string) error
}

type RoleRepo interface {
	GetRoleId(ctx context.Context, roleName string) (int32, error)
	GetRoles(ctx context.Context) ([]entity.Role, error)
}
