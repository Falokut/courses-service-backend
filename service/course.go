package service

import (
	"context"
	"courses-service/domain"
	"courses-service/entity"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Course struct {
	courseRepo CourseRepo
	fileRepo   FileRepo
	txRunner   CourseTxRunner
}

func NewCourse(courseRepo CourseRepo, txRunner CourseTxRunner, fileRepo FileRepo) Course {
	return Course{
		courseRepo: courseRepo,
		txRunner:   txRunner,
		fileRepo:   fileRepo,
	}
}

func (s Course) GetCoursesPreview(ctx context.Context, req domain.LimitOffsetRequest) ([]domain.CoursePreview, error) {
	courses, err := s.courseRepo.GetCoursesPreview(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, errors.WithMessage(err, "get courses preview")
	}
	return entityCoursesPreviewToDomain(courses), nil
}

func (s Course) GetCourse(ctx context.Context, req domain.GetCourseRequest) (*domain.Course, error) {
	course, err := s.courseRepo.GetCourse(ctx, req.CourseId)
	if err != nil {
		return nil, errors.WithMessage(err, "get course")
	}
	return &domain.Course{
		Id:        course.Id,
		AuthorFio: course.AuthorFio,
		Title:     course.Title,
		Lessons:   entityLessonsToDomain(course.Lessons),
	}, nil
}

func (s Course) GetUserCourses(ctx context.Context, userId int64) ([]domain.CoursePreview, error) {
	courses, err := s.courseRepo.GetUserCourses(ctx, userId)
	if err != nil {
		return nil, errors.WithMessage(err, "get user courses")
	}
	return entityCoursesPreviewToDomain(courses), nil
}

func (s Course) GetLectorCourses(ctx context.Context, userId int64) ([]domain.CoursePreview, error) {
	courses, err := s.courseRepo.GetCoursesByAuthorId(ctx, userId)
	if err != nil {
		return nil, errors.WithMessage(err, "get user courses")
	}
	return entityCoursesPreviewToDomain(courses), nil
}

func (s Course) Register(ctx context.Context, courseId int64, userId int64) error {
	err := s.courseRepo.Register(ctx, courseId, userId)
	if err != nil {
		return errors.WithMessage(err, "register user on course")
	}
	return nil
}

func (s Course) IsRegistered(ctx context.Context, courseId int64, userId int64) (*domain.IsRegisteredResponse, error) {
	registered, err := s.courseRepo.IsRegistered(ctx, courseId, userId)
	if err != nil {
		return nil, errors.WithMessage(err, "is registered")
	}
	return &domain.IsRegisteredResponse{IsRegistered: registered}, nil
}

func (s Course) DeleteCourse(ctx context.Context, courseId int64) error {
	err := s.courseRepo.DeleteCourse(ctx, courseId)
	if err != nil {
		return errors.WithMessage(err, "delete course")
	}
	return nil
}

func (s Course) AddCourse(ctx context.Context, req domain.AddCourseRequest) (*domain.AddCourseResponse, error) {
	previewPictureUrl := ""
	var err error
	err = s.txRunner.AddCourseTransaction(ctx, func(ctx context.Context, tx AddCourseTx) error {
		previewPictureUrl, err = s.addCourse(ctx, req, tx)
		if err != nil {
			return errors.WithMessage(err, "add course")
		}
		return nil
	})
	if err != nil {
		return nil, errors.WithMessage(err, "add course transaction")
	}
	return &domain.AddCourseResponse{
		PreviewPictureUrl: previewPictureUrl,
	}, nil
}

func (s Course) addCourse(ctx context.Context, req domain.AddCourseRequest, tx AddCourseTx) (string, error) {
	filename := uuid.NewString()
	previewPictureUrl := s.fileRepo.GetFileUrl(entity.CoursesCategory, filename)

	err := tx.AddCourse(ctx, entity.AddCourseRequest{
		AuthorId:          req.AuthorId,
		Title:             req.Title,
		PreviewPictureUrl: previewPictureUrl,
	})
	if err != nil {
		return "", errors.WithMessage(err, "add course")
	}

	err = s.fileRepo.UploadFile(ctx, entity.UploadFileReq{
		Category: entity.CoursesCategory,
		Filename: filename,
		FileBody: req.PreviewPicture,
	})
	if err != nil {
		return "", errors.WithMessage(err, "upload file")
	}
	return previewPictureUrl, nil
}

func (s Course) EditCourse(ctx context.Context, req domain.EditCourseRequest) (*domain.EditCourseResponse, error) {
	newPreviewPictureUrl := ""
	var err error
	err = s.txRunner.EditCourseTransaction(ctx, func(ctx context.Context, tx EditCourseTx) error {
		newPreviewPictureUrl, err = s.editCourse(ctx, req, tx)
		if err != nil {
			return errors.WithMessage(err, "edit course")
		}
		return nil
	})
	if err != nil {
		return nil, errors.WithMessage(err, "edit course transaction")
	}
	return &domain.EditCourseResponse{
		NewPreviewPictureUrl: newPreviewPictureUrl,
	}, nil
}

func (s Course) editCourse(ctx context.Context, req domain.EditCourseRequest, tx EditCourseTx) (string, error) {
	previewPictureUrl, err := tx.GetCoursePreviewPicture(ctx, req.CourseId)
	if err != nil {
		return "", errors.WithMessage(err, "get course preview picture")
	}
	err = s.fileRepo.DeleteFile(ctx, previewPictureUrl)
	if err != nil && !errors.Is(err, entity.ErrFileNotFound) {
		return "", errors.WithMessage(err, "delete file")
	}

	filename := uuid.NewString()
	previewPictureUrl = s.fileRepo.GetFileUrl(entity.CoursesCategory, filename)

	err = tx.EditCourse(ctx, entity.EditCourseRequest{
		Id:                req.CourseId,
		AuthorId:          req.AuthorId,
		Title:             req.Title,
		PreviewPictureUrl: previewPictureUrl,
	})
	if err != nil {
		return "", errors.WithMessage(err, "add course")
	}

	err = s.fileRepo.UploadFile(ctx, entity.UploadFileReq{
		Category: entity.CoursesCategory,
		Filename: filename,
		FileBody: req.PreviewPicture,
	})
	if err != nil {
		return "", errors.WithMessage(err, "upload file")
	}
	return previewPictureUrl, nil
}

func entityCoursesPreviewToDomain(courses []entity.CoursePreview) []domain.CoursePreview {
	domainCourses := make([]domain.CoursePreview, 0, len(courses))
	for _, course := range courses {
		domainCourses = append(domainCourses, domain.CoursePreview{
			Id:                course.Id,
			AuthorId:          course.AuthorId,
			AuthorFio:         course.AuthorFio,
			Title:             course.Title,
			PreviewPictureUrl: course.PreviewPictureUrl,
		})
	}
	return domainCourses
}

func entityLessonsToDomain(lessons []entity.Lesson) []domain.Lesson {
	domainLessons := make([]domain.Lesson, 0, len(lessons))
	for _, lesson := range lessons {
		attachments := entityLessonAttachmentsToDomain(lesson.Attachments)
		domainLessons = append(domainLessons, domain.Lesson{
			LessonNumber:  lesson.LessonNumber,
			Title:         lesson.Title,
			CreatedAt:     lesson.CreatedAt,
			LessonContent: lesson.LessonContent,
			Attachments:   attachments,
			VideoUrl:      lesson.VideoUrl,
		})
	}
	return domainLessons
}

func entityLessonAttachmentsToDomain(attachments []entity.LessonAttachment) []domain.LessonAttachment {
	domainAttachments := make([]domain.LessonAttachment, 0, len(attachments))
	for _, attachment := range attachments {
		domainAttachments = append(domainAttachments, domain.LessonAttachment{
			Id:         attachment.Id,
			LessonId:   attachment.LessonId,
			Type:       attachment.Type,
			PrettyName: attachment.PrettyName,
			Url:        attachment.Url,
		})
	}
	return domainAttachments
}
