// nolint:funlen
package routes

import (
	"courses-service/controller"
	"net/http"

	"github.com/Falokut/go-kit/http/endpoint"
	"github.com/Falokut/go-kit/http/router"
)

type Router struct {
	Auth controller.Auth
}

func (r Router) InitRoutes(authMiddleware AuthMiddleware, wrapper endpoint.Wrapper) *router.Router {
	mux := router.New()
	for _, desc := range endpointDescriptors(r) {
		endpointWrapper := wrapper
		switch {
		case desc.IsAdmin:
			endpointWrapper = wrapper.WithMiddlewares(authMiddleware.AdminAuthToken())
		case desc.IsTeacher:
			endpointWrapper = wrapper.WithMiddlewares(authMiddleware.TeacherAuthToken())
		case desc.NeedUserAuth:
			endpointWrapper = wrapper.WithMiddlewares(authMiddleware.UserAuthToken())
		}
		mux.Handler(desc.Method, desc.Path, endpointWrapper.Endpoint(desc.Handler))
	}

	return mux
}

type EndpointDescriptor struct {
	Method       string
	Path         string
	IsAdmin      bool
	IsTeacher    bool
	NeedUserAuth bool
	Handler      any
}

func endpointDescriptors(r Router) []EndpointDescriptor {
	return []EndpointDescriptor{
		{
			Method:  http.MethodPost,
			Path:    "/auth/login",
			Handler: r.Auth.Login,
		}, {
			Method:  http.MethodPost,
			Path:    "/auth/register",
			Handler: r.Auth.Register,
		}, {
			Method:  http.MethodGet,
			Path:    "/auth/refresh_access_token",
			Handler: r.Auth.RefreshAccessToken,
		}, {
			Method:  http.MethodGet,
			Path:    "/auth/get_user_role",
			Handler: r.Auth.GetRole,
		},
	}
}
