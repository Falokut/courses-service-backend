package routes

import (
	"context"
	"courses-service/domain"
	"courses-service/entity"
	"net/http"
	"slices"

	"github.com/Falokut/go-kit/jwt"

	http2 "github.com/Falokut/go-kit/http"
	"github.com/Falokut/go-kit/http/apierrors"
	"github.com/Falokut/go-kit/http/types"
)

type AuthMiddleware struct {
	accessTokenSecret string
}

func NewAuthMiddleware(accessTokenSecret string) AuthMiddleware {
	return AuthMiddleware{
		accessTokenSecret: accessTokenSecret,
	}
}

func (m AuthMiddleware) AdminAuthToken() http2.Middleware {
	return AuthToken(m.accessTokenSecret, domain.AdminType)
}

func (m AuthMiddleware) UserAuthToken() http2.Middleware {
	return AuthToken(m.accessTokenSecret, domain.StudentType, domain.AdminType)
}

func (m AuthMiddleware) TeacherAuthToken() http2.Middleware {
	return AuthToken(m.accessTokenSecret, domain.LectorType, domain.AdminType)
}

func AuthToken(tokenSecret string, roles ...string) http2.Middleware {
	return func(next http2.HandlerFunc) http2.HandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			token := types.BearerToken{}
			err := token.FromRequestHeader(r)
			if err != nil {
				return err
			}

			userInfo := entity.TokenUserInfo{}
			err = jwt.ParseToken(token.Token, tokenSecret, &userInfo)
			if err != nil {
				return err
			}
			domain.UserIdToContext(ctx, userInfo.UserId)
			if len(roles) == 0 {
				return next(ctx, w, r)
			}
			if !slices.Contains(roles, userInfo.RoleName) {
				return apierrors.NewForbiddenError("access forbidden")
			}
			return next(ctx, w, r)
		}
	}
}
