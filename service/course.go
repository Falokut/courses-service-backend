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
