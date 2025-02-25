// nolint:funlen
package routes

import (
	"courses-service/controller"
	"courses-service/domain"
	"net/http"

	"github.com/Falokut/go-kit/http/endpoint"
	"github.com/Falokut/go-kit/http/router"
)

type Router struct {
	Auth controller.Auth
	User controller.User
	Role controller.Role
}

func (r Router) InitRoutes(authMiddleware AuthMiddleware, wrapper endpoint.Wrapper) *router.Router {
	mux := router.New()
	disableCors := DisableCors{}
	mux.InternalRouter().GlobalOPTIONS = disableCors
	for _, desc := range endpointDescriptors(r) {
		endpointWrapper := wrapper.WithMiddlewares(disableCors.Middleware)
		if len(desc.AllowedRoles) > 0 {
			endpointWrapper = endpointWrapper.WithMiddlewares(authMiddleware.AuthRole(desc.AllowedRoles...))
		}
		mux.Handler(desc.Method, desc.Path, endpointWrapper.Endpoint(desc.Handler))
	}

	return mux
}

type EndpointDescriptor struct {
	Method       string
	Path         string
	AllowedRoles []string
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
			Path:    "/auth/logout",
			Handler: r.Auth.Logout,
		}, {
			Method:       http.MethodPost,
			Path:         "/auth/register",
			Handler:      r.Auth.Register,
			AllowedRoles: []string{domain.AdminType},
		}, {
			Method:  http.MethodGet,
			Path:    "/users/get_role",
			Handler: r.User.GetRole,
		}, {
			Method:       http.MethodGet,
			Path:         "/users",
			Handler:      r.User.GetUsers,
			AllowedRoles: []string{domain.AdminType},
		}, {
			Method:       http.MethodDelete,
			Path:         "/users/:userId",
			Handler:      r.User.DeleteUser,
			AllowedRoles: []string{domain.AdminType},
		}, {
			Method:  http.MethodGet,
			Path:    "/roles",
			Handler: r.Role.GetRoles,
		},
	}
}
