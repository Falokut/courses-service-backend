package service

import (
	"context"
	"courses-service/conf"
	"courses-service/domain"
	"courses-service/entity"
	"time"

	"github.com/google/uuid"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	Register(ctx context.Context, req entity.RegisterUser) error
}

type AuthTxRunner interface {
	LoginTransaction(ctx context.Context, txFunc func(ctx context.Context, tx LoginTx) error) error
}

type LoginTx interface {
	GetUserByUsername(ctx context.Context, username string) (entity.User, error)
	InsertSession(ctx context.Context, session entity.Session) error
}

type Auth struct {
	cfg      conf.Auth
	userRepo UserRepo
	txRunner AuthTxRunner
}

func NewAuth(cfg conf.Auth, userRepo UserRepo, txRunner AuthTxRunner) Auth {
	return Auth{
		cfg:      cfg,
		userRepo: userRepo,
		txRunner: txRunner,
	}
}

func (s Auth) Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, error) {
	var sessionId string
	var err error

	err = s.txRunner.LoginTransaction(ctx,
		func(ctx context.Context, tx LoginTx) error {
			sessionId, err = s.login(ctx, req, tx)
			if err != nil {
				return errors.WithMessage(err, "login")
			}
			return nil
		})
	if err != nil {
		return nil, errors.WithMessage(err, "login transaction")
	}

	return &domain.LoginResponse{
		SessionId: sessionId,
	}, nil
}

func (s Auth) login(ctx context.Context, req domain.LoginRequest, tx LoginTx) (string, error) {
	user, err := tx.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return "", errors.WithMessage(err, "get user by username")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	session := entity.Session{
		Id:        uuid.NewString(),
		UserId:    user.Id,
		CreatedAt: time.Now().UTC(),
	}
	err = tx.InsertSession(ctx, session)
	if err != nil {
		return "", errors.WithMessage(err, "insert session")
	}
	return session.Id, nil
}

func (s Auth) Register(ctx context.Context, req domain.RegisterRequest) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), s.cfg.BcryptCost)
	if err != nil {
		return errors.WithMessage(err, "generate from passport")
	}

	err = s.userRepo.Register(ctx, entity.RegisterUser{
		Username:     req.Username,
		Fio:          req.Fio,
		PasswordHash: string(passwordHash),
		RoleId:       req.RoleId,
	})
	if err != nil {
		return errors.WithMessage(err, "register")
	}
	return nil
}
