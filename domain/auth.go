package domain

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

type userIdKey struct{}

func UserIdToContext(ctx context.Context, userId int64) context.Context {
	return context.WithValue(ctx, userIdKey{}, fmt.Sprint(userId))
}

func UserIdFromContext(ctx context.Context) int64 {
	userId, err := strconv.Atoi(fmt.Sprint(ctx.Value(userIdKey{})))
	if err != nil {
		return -1
	}
	return int64(userId)
}

const (
	AdminType   = "admin"
	StudentType = "student"
	LectorType  = "lector"
)

type LoginRequest struct {
	Username string `validate:"min=4,max=50"`
	Password string `validate:"min=6,max=20"`
}

type LoginResponse struct {
	SessionId string
}

type RegisterRequest struct {
	Username string `validate:"min=4,max=50"`
	Fio      string `validate:"min=6,max=60"`
	Password string `validate:"min=6,max=20"`
	RoleId   int32  `validate:"required"`
}

type TerminateSessionRequest struct {
	SessionId string `validate:"required"`
}

type Session struct {
	Id        string
	CreatedAt time.Time
}
