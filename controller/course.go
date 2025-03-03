package controller

import (
	"context"
	"courses-service/domain"
)

type CourseService interface {
	GetCoursesPreview(ctx context.Context, req domain.LimitOffsetRequest) ([]domain.CoursePreview, error)
	GetCourse(ctx context.Context, req domain.GetCourseRequest) (*domain.Course, error)
	GetUserCourses(ctx context.Context, userId int64) ([]domain.CoursePreview, error)
	Register(ctx context.Context, courseId int64, userId int64) error
	IsRegistered(ctx context.Context, courseId int64, userId int64) (domain.IsRegisteredResponse, error)
}

type Course struct {
	service CourseService
}

func NewCourse(service CourseService) Course {
	return Course{
		service: service,
	}
}

// GetCoursesPreview
//
//	@Tags		course
//	@Summary	Получить список курсов
//	@Accept		json
//	@Produce	json
//
//
//	@Success	200	{array}		domain.CoursePreview
//	@Failure	400	{object}	apierrors.Error
//	@Failure	500	{object}	apierrors.Error
//	@Router		/courses [GET]
func (c Course) GetCoursesPreview(ctx context.Context, req domain.LimitOffsetRequest) ([]domain.CoursePreview, error) {
	return c.service.GetCoursesPreview(ctx, req)
}

// GetCourse
//
//	@Tags		course
//	@Summary	Получить курс
//	@Accept		json
//	@Produce	json
//
//	@Param		courseId	query		int	true	"идентификатор курса"
//
//	@Success	200			{object}	domain.Course
//	@Failure	400			{object}	apierrors.Error
//	@Failure	404			{object}	apierrors.Error
//	@Failure	500			{object}	apierrors.Error
//	@Router		/courses/by_id [GET]
func (c Course) GetCourse(ctx context.Context, req domain.GetCourseRequest) (*domain.Course, error) {
	return c.service.GetCourse(ctx, req)
}

// Register
//
//	@Tags		course
//	@Summary	Зарегистрироваться на курс
//	@Produce	json
//
//	@Param		body	body		domain.CourseIdRequest	true	"тело запроса"
//
//	@Success	200		{object}	any
//	@Failure	400		{object}	apierrors.Error
//	@Failure	401		{object}	apierrors.Error
//	@Failure	403		{object}	apierrors.Error
//	@Failure	500		{object}	apierrors.Error
//
//	@Security	Bearer
//
//	@Router		/courses/register [POST]
func (c Course) Register(ctx context.Context, req domain.CourseIdRequest) error {
	return c.service.Register(ctx, req.CourseId, domain.UserIdFromContext(ctx))
}

// GetUserCourses
//
//	@Tags		course
//	@Summary	Получить курсы пользователя
//	@Produce	json
//
//	@Success	200	{array}		domain.CoursePreview
//	@Failure	400	{object}	apierrors.Error
//	@Failure	401	{object}	apierrors.Error
//	@Failure	403	{object}	apierrors.Error
//	@Failure	500	{object}	apierrors.Error
//
//	@Security	Bearer
//
//	@Router		/courses/user_courses [GET]
func (c Course) GetUserCourses(ctx context.Context) ([]domain.CoursePreview, error) {
	return c.service.GetUserCourses(ctx, domain.UserIdFromContext(ctx))
}

// IsRegistered
//
//	@Tags		course
//	@Summary	Проверить зарегистрирован пользователь на курс
//	@Produce	json
//
//	@Param		courseId	query		int	true	"идентификатор курса"
//	@Success	200			{object}	domain.IsRegisteredResponse
//	@Failure	400			{object}	apierrors.Error
//	@Failure	401			{object}	apierrors.Error
//	@Failure	403			{object}	apierrors.Error
//	@Failure	500			{object}	apierrors.Error
//
//	@Security	Bearer
//
//	@Router		/courses/is_registered [GET]
func (c Course) IsRegistered(ctx context.Context, req domain.CourseIdRequest) (domain.IsRegisteredResponse, error) {
	return c.service.IsRegistered(ctx, req.CourseId, domain.UserIdFromContext(ctx))
}
