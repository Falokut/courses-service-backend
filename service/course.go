package service

import (
	"context"
	"courses-service/domain"

	"github.com/pkg/errors"
)

type Course struct {
	repo CourseRepo
}

func NewCourse(repo CourseRepo) Course {
	return Course{
		repo: repo,
	}
}

func (s Course) GetCoursesPreview(ctx context.Context, req domain.LimitOffsetRequest) ([]domain.CoursePreview, error) {
	courses, err := s.repo.GetCoursesPreview(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, errors.WithMessage(err, "get courses preview")
	}

	domainCourses := make([]domain.CoursePreview, 0, len(courses))
	for _, course := range courses {
		domainCourses = append(domainCourses, domain.CoursePreview{
			Id:                course.Id,
			AuthorFio:         course.AuthorFio,
			Title:             course.Title,
			PreviewPictureUrl: course.PreviewPictureUrl,
		})
	}
	return domainCourses, nil
}

func (s Course) GetCourse(ctx context.Context, req domain.GetCourseRequest) (*domain.Course, error) {
	course, err := s.repo.GetCourse(ctx, req.CourseId)
	if err != nil {
		return nil, errors.WithMessage(err, "get course")
	}

	courseLessons := make([]domain.Lesson, 0, len(course.Lessons))
	for _, lesson := range course.Lessons {
		courseLessons = append(courseLessons, domain.Lesson{
			LessonNumber:  lesson.LessonNumber,
			Title:         lesson.Title,
			CreatedAt:     lesson.CreatedAt,
			LessonContent: lesson.LessonContent,
		})
	}
	return &domain.Course{
		Id:        course.Id,
		AuthorFio: course.AuthorFio,
		Title:     course.Title,
		Lessons:   courseLessons,
	}, nil
}

func (s Course) GetUserCourses(ctx context.Context, userId int64) ([]domain.CoursePreview, error) {
	courses, err := s.repo.GetUserCourses(ctx, userId)
	if err != nil {
		return nil, errors.WithMessage(err, "get user courses")
	}

	domainCourses := make([]domain.CoursePreview, 0, len(courses))
	for _, course := range courses {
		domainCourses = append(domainCourses, domain.CoursePreview{
			Id:                course.Id,
			AuthorFio:         course.AuthorFio,
			Title:             course.Title,
			PreviewPictureUrl: course.PreviewPictureUrl,
		})
	}
	return domainCourses, nil
}

func (s Course) Register(ctx context.Context, courseId int64, userId int64) error {
	err := s.repo.Register(ctx, courseId, userId)
	if err != nil {
		return errors.WithMessage(err, "get user courses")
	}
	return nil
}

func (s Course) IsRegistered(ctx context.Context, courseId int64, userId int64) (domain.IsRegisteredResponse, error) {
	registered, err := s.repo.IsRegistered(ctx, courseId, userId)
	if err != nil {
		return domain.IsRegisteredResponse{}, errors.WithMessage(err, "is registered")
	}
	return domain.IsRegisteredResponse{IsRegistered: registered}, nil
}
