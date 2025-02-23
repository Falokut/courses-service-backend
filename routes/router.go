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
	User controller.User
}

func (r Router) InitRoutes(authMiddleware AuthMiddleware, wrapper endpoint.Wrapper) *router.Router {
	mux := router.New()
	disableCors := DisableCors{}
	mux.InternalRouter().GlobalOPTIONS = disableCors

	for _, desc := range endpointDescriptors(r) {
		endpointWrapper := wrapper.WithMiddlewares(disableCors.Middleware)
		switch {
		case desc.IsAdmin:
			endpointWrapper = wrapper.WithMiddlewares(authMiddleware.AdminAuthToken())
		case desc.IsLector:
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
	IsLector     bool
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
			Path:    "/user/get_role",
			Handler: r.User.GetRole,
		},
	}
}
