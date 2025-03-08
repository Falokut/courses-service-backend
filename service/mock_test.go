// nolint:wrapcheck
package service_test

import (
	"context"
	"courses-service/entity"

	"github.com/stretchr/testify/mock"
)

type MockLessonRepo struct {
	mock.Mock
}

func (m *MockLessonRepo) CheckLessonOwnership(ctx context.Context, userId int64, lessonId int64) (bool, error) {
	args := m.Called(ctx, userId, lessonId)
	return args.Bool(0), args.Error(1)
}

func (m *MockLessonRepo) CreateLesson(ctx context.Context, lesson entity.CreateLessonRequest) error {
	args := m.Called(ctx, lesson)
	return args.Error(0)
}

func (m *MockLessonRepo) EditTitle(ctx context.Context, id int64, newTitle string) error {
	args := m.Called(ctx, id, newTitle)
	return args.Error(0)
}

func (m *MockLessonRepo) EditContent(ctx context.Context, id int64, content string) error {
	args := m.Called(ctx, id, content)
	return args.Error(0)
}

func (m *MockLessonRepo) UnattachFile(ctx context.Context, attachmentId int64) error {
	args := m.Called(ctx, attachmentId)
	return args.Error(0)
}

type MockCourseRepo struct {
	mock.Mock
}

func (m *MockCourseRepo) CheckCourseOwnership(ctx context.Context, userId int64, courseId int64) (bool, error) {
	args := m.Called(ctx, userId, courseId)
	return args.Bool(0), args.Error(1)
}
