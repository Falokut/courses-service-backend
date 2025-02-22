package domain

import (
	"context"
	"strconv"

	"github.com/Falokut/go-kit/jwt"
)

type userIdKey struct{}

func UserIdToContext(ctx context.Context, userId int64) context.Context {
	return context.WithValue(ctx, userIdKey{}, userId)
}

func UserIdFromContext(ctx context.Context) int64 {
	userId, err := strconv.Atoi(ctx.Value(userIdKey{}).(string))
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
	AccessToken  jwt.TokenResponse
	RefreshToken jwt.TokenResponse
}

type RegisterRequest struct {
	Username string `validate:"min=4,max=50"`
	Fio      string `validate:"min=6,max=60"`
	Password string `validate:"min=6,max=20"`
	RoleId   int32  `validate:"required"`
}

type GetRoleResponse struct {
	RoleName string
}
