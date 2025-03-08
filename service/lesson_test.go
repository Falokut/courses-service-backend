package service_test

import (
	"context"
	"testing"

	"courses-service/domain"
	"courses-service/service"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestLessonServiceTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(LessonServiceTestSuite))
}

type LessonServiceTestSuite struct {
	suite.Suite
	lessonRepo *MockLessonRepo
	courseRepo *MockCourseRepo
	lesson     service.Lesson
}

func (suite *LessonServiceTestSuite) SetupTest() {
	suite.lessonRepo = &MockLessonRepo{}
	suite.courseRepo = &MockCourseRepo{}
	suite.lesson = service.NewLesson(suite.lessonRepo, suite.courseRepo, nil, nil)
}

func (s *LessonServiceTestSuite) TestCreateLesson_Success() {
	ctx := context.Background()
	req := domain.CreateLessonRequest{
		CourseId:     1,
		LessonNumber: 1,
		Title:        "Test Lesson",
	}

	s.courseRepo.On("CheckCourseOwnership", ctx, mock.Anything, req.CourseId).Return(true, nil)
	s.lessonRepo.On("CreateLesson", ctx, mock.Anything).Return(nil)

	err := s.lesson.CreateLesson(ctx, req)
	s.Require().NoError(err)
	s.courseRepo.AssertExpectations(s.T())
	s.lessonRepo.AssertExpectations(s.T())
}

func (s *LessonServiceTestSuite) TestEditTitle_Success() {
	ctx := context.Background()
	req := domain.EditLessonTitleRequest{
		LessonId: 1,
		NewTitle: "Updated Title",
	}

	s.lessonRepo.On("CheckLessonOwnership", ctx, mock.Anything, req.LessonId).Return(true, nil)
	s.lessonRepo.On("EditTitle", ctx, req.LessonId, req.NewTitle).Return(nil)

	err := s.lesson.EditTitle(ctx, req)
	s.Require().NoError(err)
	s.lessonRepo.AssertExpectations(s.T())
}
