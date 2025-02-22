package entity

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

type TokenUserInfo struct {
	UserId   int64  `json:"user_id"`
	RoleName string `json:"role"`
}

func (i *TokenUserInfo) FromMap(m map[string]any) error {
	userId, err := strconv.Atoi(m["userId"].(string))
	if err != nil {
		return errors.WithMessage(err, "parse user id")
	}
	i.UserId = int64(userId)
	i.RoleName = m["roleName"].(string)
	return nil
}

func (i *TokenUserInfo) ToMap() map[string]any {
	return map[string]any{
		"userId":   fmt.Sprint(i.UserId),
		"roleName": i.RoleName,
	}
}
