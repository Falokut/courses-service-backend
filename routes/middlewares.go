package routes

import (
	"context"
	"courses-service/domain"
	"courses-service/entity"
	"net/http"
	"slices"

	"github.com/pkg/errors"

	http2 "github.com/Falokut/go-kit/http"
	"github.com/Falokut/go-kit/http/apierrors"
	"github.com/Falokut/go-kit/http/types"
)

type AuthRepo interface {
	GetUserSession(ctx context.Context, sessionId string) (entity.UserSession, error)
}

type AuthMiddleware struct {
	authRepo AuthRepo
}

func NewAuthMiddleware(authRepo AuthRepo) AuthMiddleware {
	return AuthMiddleware{
		authRepo: authRepo,
	}
}

func (m AuthMiddleware) AuthRole(roles ...string) http2.Middleware {
	return func(next http2.HandlerFunc) http2.HandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			token := types.BearerToken{}
			err := token.FromRequestHeader(r)
			if err != nil {
				return apierrors.NewUnauthorizedError(domain.ErrInvalidAuthData.Error())
			}
			userSession, err := m.authRepo.GetUserSession(ctx, token.Token)
			switch {
			case errors.Is(err, domain.ErrSessionNotFound):
				return apierrors.NewForbiddenError(domain.ErrForbidden.Error())
			case err != nil:
				return errors.WithMessage(err, "get user session")
			}

			ctx = domain.UserIdToContext(ctx, userSession.UserId)
			if len(roles) == 0 {
				return next(ctx, w, r)
			}
			if !slices.Contains(roles, userSession.RoleName) {
				return apierrors.NewForbiddenError(domain.ErrForbidden.Error())
			}
			return next(ctx, w, r)
		}
	}
}

type DisableCors struct{}

func (c DisableCors) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "*")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Methods", "*")
	w.Header().Add("Access-Control-Max-Age", "1728000")
	w.Header().Add("Access-Control-Allow-Credentials", "false")
}

func (c DisableCors) Middleware(next http2.HandlerFunc) http2.HandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		c.ServeHTTP(w, r)
		return next(ctx, w, r)
	}
}
