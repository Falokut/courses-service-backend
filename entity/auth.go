package entity

import (
	"time"
)

type Session struct {
	Id        string
	UserId    int64
	CreatedAt time.Time
}

type UserSession struct {
	Id        string
	UserId    int64
	CreatedAt time.Time
	RoleName  string
}
