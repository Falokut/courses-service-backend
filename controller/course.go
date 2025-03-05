package controller

import (
	"context"
	"courses-service/domain"
	"errors"
	"net/http"

	"github.com/Falokut/go-kit/http/apierrors"
)

type CourseService interface {
	GetCoursesPreview(ctx context.Context, req domain.LimitOffsetRequest) ([]domain.CoursePreview, error)
	GetCourse(ctx context.Context, req domain.GetCourseRequest) (*domain.Course, error)
	GetUserCourses(ctx context.Context, userId int64) ([]domain.CoursePreview, error)
	GetLectorCourses(ctx context.Context, userId int64) ([]domain.CoursePreview, error)
	Register(ctx context.Context, courseId int64, userId int64) error
	IsRegistered(ctx context.Context, courseId int64, userId int64) (*domain.IsRegisteredResponse, error)
	DeleteCourse(ctx context.Context, courseId int64) error
	AddCourse(ctx context.Context, req domain.AddCourseRequest) (*domain.AddCourseResponse, error)
	EditCourse(ctx context.Context, req domain.EditCourseRequest) (*domain.EditCourseResponse, error)
	ReorderLessons(ctx context.Context, req domain.EditCourseLessonsOrderingRequest) error
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

// GetLectorCourses
//
//	@Tags		course
//	@Summary	Получить курсы лектора (получение курсов, автором которых является пользователь)
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
//	@Router		/courses/lector_courses [GET]
func (c Course) GetLectorCourses(ctx context.Context) ([]domain.CoursePreview, error) {
	return c.service.GetLectorCourses(ctx, domain.UserIdFromContext(ctx))
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
func (c Course) IsRegistered(ctx context.Context, req domain.CourseIdRequest) (*domain.IsRegisteredResponse, error) {
	return c.service.IsRegistered(ctx, req.CourseId, domain.UserIdFromContext(ctx))
}

// DeleteCourse
//
//	@Tags		course
//	@Summary	Удалить курс
//	@Produce	json
//
//	@Param		courseId	query		int	true	"идентификатор курса"
//
//	@Success	200			{object}	any
//	@Failure	400			{object}	apierrors.Error
//	@Failure	401			{object}	apierrors.Error
//	@Failure	403			{object}	apierrors.Error
//	@Failure	500			{object}	apierrors.Error
//
//	@Security	Bearer
//
//	@Router		/courses [DELETE]
func (c Course) DeleteCourse(ctx context.Context, req domain.CourseIdRequest) error {
	return c.service.DeleteCourse(ctx, req.CourseId)
}

// AddCourse
//
//	@Tags		course
//	@Summary	Создать курс
//	@Accept		json
//	@Produce	json
//
//	@Param		body	body		domain.AddCourseRequest	true	"тело запроса"
//
//	@Success	200		{object}	domain.AddCourseResponse
//	@Failure	400		{object}	apierrors.Error
//	@Failure	401		{object}	apierrors.Error
//	@Failure	403		{object}	apierrors.Error
//	@Failure	500		{object}	apierrors.Error
//
//	@Security	Bearer
//
//	@Router		/courses [POST]
func (c Course) AddCourse(ctx context.Context, req domain.AddCourseRequest) (*domain.AddCourseResponse, error) {
	return c.service.AddCourse(ctx, req)
}

// EditCourse
//
//	@Tags		course
//	@Summary	Редактировать курс
//	@Accept		json
//	@Produce	json
//
//	@Param		body	body		domain.EditCourseRequest	true	"тело запроса"
//
//	@Success	200		{object}	domain.EditCourseResponse
//	@Failure	400		{object}	apierrors.Error
//	@Failure	401		{object}	apierrors.Error
//	@Failure	403		{object}	apierrors.Error
//	@Failure	404		{object}	apierrors.Error
//	@Failure	500		{object}	apierrors.Error
//
//	@Security	Bearer
//
//	@Router		/courses/edit [POST]
func (c Course) EditCourse(ctx context.Context, req domain.EditCourseRequest) (*domain.EditCourseResponse, error) {
	resp, err := c.service.EditCourse(ctx, req)
	switch {
	case errors.Is(err, domain.ErrCourseNotFound):
		return nil, apierrors.New(http.StatusNotFound, domain.ErrCodeCourseNotFound, domain.ErrCourseNotFound.Error(), err)
	default:
		return resp, err
	}
}

// ReorderLessons
//
//	@Tags		course
//	@Summary	Изменить порядок уроков в курсе
//	@Accept		json
//	@Produce	json
//
//	@Param		body	body		domain.EditCourseLessonsOrderingRequest	true	"тело запроса"
//
//	@Success	200		{object}	any
//	@Failure	400		{object}	apierrors.Error
//	@Failure	401		{object}	apierrors.Error
//	@Failure	403		{object}	apierrors.Error
//	@Failure	404		{object}	apierrors.Error
//	@Failure	500		{object}	apierrors.Error
//
//	@Security	Bearer
//
//	@Router		/courses/reorder_lessons [POST]
func (c Course) ReorderLessons(ctx context.Context, req domain.EditCourseLessonsOrderingRequest) error {
	err := c.service.ReorderLessons(ctx, req)
	switch {
	case errors.Is(err, domain.ErrCourseNotFound):
		return apierrors.New(http.StatusNotFound, domain.ErrCodeCourseNotFound, domain.ErrCourseNotFound.Error(), err)
	default:
		return err
	}
}
