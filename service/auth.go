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

type Auth struct {
	cfg      conf.Auth
	authRepo AuthRepo
	userRepo UserRepo
	roleRepo RoleRepo
	txRunner AuthTxRunner
}

func NewAuth(
	cfg conf.Auth,
	authRepo AuthRepo,
	userRepo UserRepo,
	roleRepo RoleRepo,
	txRunner AuthTxRunner,
) Auth {
	return Auth{
		cfg:      cfg,
		authRepo: authRepo,
		userRepo: userRepo,
		roleRepo: roleRepo,
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

func (s Auth) InitAdmin(ctx context.Context, adminAuthData conf.InitAdmin) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(adminAuthData.Password), s.cfg.BcryptCost)
	if err != nil {
		return errors.WithMessage(err, "generate from passport")
	}

	adminRoleId, err := s.roleRepo.GetRoleId(ctx, domain.AdminType)
	if err != nil {
		return errors.WithMessage(err, "get admin role id")
	}
	err = s.userRepo.UpsertUser(ctx, entity.UpsertUser{
		Fio:          "ADMIN ADMIN",
		Username:     adminAuthData.Username,
		PasswordHash: string(passwordHash),
		RoleId:       adminRoleId,
	})
	if err != nil {
		return errors.WithMessage(err, "upsert user")
	}
	return nil
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

func (s Auth) Logout(ctx context.Context, sessionId string) error {
	err := s.authRepo.DeleteSession(ctx, sessionId)
	if err != nil {
		return errors.WithMessage(err, "delete session")
	}
	return nil
}
