package controller

import (
	"context"
	"courses-service/domain"
)

type CourseService interface {
	GetCoursesPreview(ctx context.Context, req domain.LimitOffsetRequest) ([]domain.CoursePreview, error)
	GetCourse(ctx context.Context, req domain.GetCourseRequest) (*domain.Course, error)
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
//
//	@Success	200	{object}	domain.Course
//	@Failure	400	{object}	apierrors.Error
//	@Failure	404	{object}	apierrors.Error
//	@Failure	500	{object}	apierrors.Error
//	@Router		/courses/{courseId} [GET]
func (c Course) GetCourse(ctx context.Context, req domain.GetCourseRequest) (*domain.Course, error) {
	return c.service.GetCourse(ctx, req)
}
