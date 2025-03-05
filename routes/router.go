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
	Auth   controller.Auth
	User   controller.User
	Role   controller.Role
	Course controller.Course
	Lesson controller.Lesson
}

func (r Router) InitRoutes(
	authMiddleware AuthMiddleware,
	wrapper endpoint.Wrapper,
	wrapperWithoutMaxBodySize endpoint.Wrapper,
) *router.Router {
	mux := router.New()
	disableCors := DisableCors{}
	mux.InternalRouter().GlobalOPTIONS = disableCors

	for _, desc := range endpointDescriptors(r) {
		var endpointWrapper endpoint.Wrapper
		if desc.DisableMaxBodySize {
			endpointWrapper = wrapperWithoutMaxBodySize.WithMiddlewares(disableCors.Middleware)
		} else {
			endpointWrapper = wrapper.WithMiddlewares(disableCors.Middleware)
		}
		if len(desc.AllowedRoles) > 0 {
			endpointWrapper = endpointWrapper.WithMiddlewares(authMiddleware.AuthRole(desc.AllowedRoles...))
		}
		mux.Handler(desc.Method, desc.Path, endpointWrapper.Endpoint(desc.Handler))
	}

	return mux
}

type EndpointDescriptor struct {
	Method             string
	Path               string
	AllowedRoles       []string
	Handler            any
	DisableMaxBodySize bool
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
			Path:         "/auth/terminate_session",
			Handler:      r.Auth.TerminateSession,
			AllowedRoles: []string{domain.AdminType, domain.StudentType, domain.LectorType},
		}, {
			Method:       http.MethodGet,
			Path:         "/auth/sessions",
			Handler:      r.Auth.SessionsList,
			AllowedRoles: []string{domain.AdminType, domain.StudentType, domain.LectorType},
		}, {
			Method:       http.MethodPost,
			Path:         "/auth/register",
			Handler:      r.Auth.Register,
			AllowedRoles: []string{domain.AdminType},
		},

		{
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
			Path:         "/users",
			Handler:      r.User.DeleteUser,
			AllowedRoles: []string{domain.AdminType},
		}, {
			Method:       http.MethodPost,
			Path:         "/users",
			Handler:      r.User.EditUser,
			AllowedRoles: []string{domain.AdminType},
		}, {
			Method:       http.MethodGet,
			Path:         "/users/profile",
			Handler:      r.User.GetUserProfile,
			AllowedRoles: []string{domain.AdminType, domain.StudentType, domain.LectorType},
		}, {
			Method:  http.MethodGet,
			Path:    "/users/lectors",
			Handler: r.User.GetLectors,
		}, {
			Method:  http.MethodGet,
			Path:    "/roles",
			Handler: r.Role.GetRoles,
		},

		{
			Method:  http.MethodGet,
			Path:    "/courses",
			Handler: r.Course.GetCoursesPreview,
		}, {
			Method:       http.MethodDelete,
			Path:         "/courses",
			Handler:      r.Course.DeleteCourse,
			AllowedRoles: []string{domain.AdminType},
		}, {
			Method:             http.MethodPost,
			Path:               "/courses",
			Handler:            r.Course.AddCourse,
			AllowedRoles:       []string{domain.AdminType},
			DisableMaxBodySize: true,
		}, {
			Method:             http.MethodPost,
			Path:               "/courses/edit",
			Handler:            r.Course.EditCourse,
			AllowedRoles:       []string{domain.AdminType},
			DisableMaxBodySize: true,
		}, {
			Method:  http.MethodGet,
			Path:    "/courses/by_id",
			Handler: r.Course.GetCourse,
		}, {
			Method:       http.MethodPost,
			Path:         "/courses/register",
			Handler:      r.Course.Register,
			AllowedRoles: []string{domain.AdminType, domain.StudentType, domain.LectorType},
		}, {
			Method:       http.MethodGet,
			Path:         "/courses/is_registered",
			Handler:      r.Course.IsRegistered,
			AllowedRoles: []string{domain.AdminType, domain.StudentType, domain.LectorType},
		}, {
			Method:       http.MethodGet,
			Path:         "/courses/user_courses",
			Handler:      r.Course.GetUserCourses,
			AllowedRoles: []string{domain.AdminType, domain.StudentType},
		}, {
			Method:       http.MethodGet,
			Path:         "/courses/lector_courses",
			Handler:      r.Course.GetLectorCourses,
			AllowedRoles: []string{domain.LectorType},
		}, {
			Method:       http.MethodGet,
			Path:         "/courses/reorder_lessons",
			Handler:      r.Course.ReorderLessons,
			AllowedRoles: []string{domain.LectorType},
		}, {
			Method:  http.MethodGet,
			Path:    "/courses/stats",
			Handler: r.Course.Stats,
		},

		{
			Method:       http.MethodPost,
			Path:         "/lessons",
			Handler:      r.Lesson.CreateLesson,
			AllowedRoles: []string{domain.LectorType},
		}, {
			Method:       http.MethodPost,
			Path:         "/lessons/edit_title",
			Handler:      r.Lesson.EditTitle,
			AllowedRoles: []string{domain.LectorType},
		}, {
			Method:       http.MethodPost,
			Path:         "/lessons/edit_content",
			Handler:      r.Lesson.EditLessonContent,
			AllowedRoles: []string{domain.LectorType},
		}, {
			Method:             http.MethodPost,
			Path:               "/lessons/add_video",
			Handler:            r.Lesson.AddLessonVideo,
			AllowedRoles:       []string{domain.LectorType},
			DisableMaxBodySize: true,
		}, {
			Method:             http.MethodPost,
			Path:               "/lessons/attach_file",
			Handler:            r.Lesson.AttachFileToLesson,
			AllowedRoles:       []string{domain.LectorType},
			DisableMaxBodySize: true,
		}, {
			Method:       http.MethodDelete,
			Path:         "/lessons/delete_video",
			Handler:      r.Lesson.DeleteLessonVideo,
			AllowedRoles: []string{domain.LectorType},
		}, {
			Method:       http.MethodDelete,
			Path:         "/lessons/unattach_file",
			Handler:      r.Lesson.UnattachFileFromLesson,
			AllowedRoles: []string{domain.LectorType},
		},
	}
}
