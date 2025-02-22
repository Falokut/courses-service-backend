package domain

import (
	"github.com/pkg/errors"
)

var (
	ErrUserAlreadyExists      = errors.New("пользователь уже существует")
	ErrUserNotFound           = errors.New("пользователь не найден")
	ErrInvalidCredentials     = errors.New("невалидные данные для авторизации")
	ErrUserOperationForbidden = errors.New("данная операция запрещена для пользователя")
	ErrWrongSecret            = errors.New("неверный пароль")
	ErrForbidden              = errors.New("доступ запрещён")
)

const (
	ErrCodeInvalidArgument = 400

	ErrCodeUserNotFound      = 604
	ErrCodeUserAlreadyExists = 605
	ErrCodeWrongSecret       = 606
)
