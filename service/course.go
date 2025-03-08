package service

import (
	"bytes"
	"context"
	"courses-service/domain"
	"courses-service/entity"
	"encoding/csv"
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type CourseRepo interface {
	GetCoursesPreview(ctx context.Context, limit int32, offset int32) ([]entity.CoursePreview, error)
	GetCourse(ctx context.Context, id int64) (entity.Course, error)
	GetUserCourses(ctx context.Context, userId int64) ([]entity.CoursePreview, error)
	GetCoursesByAuthorId(ctx context.Context, userId int64) ([]entity.CoursePreview, error)
	Register(ctx context.Context, courseId int64, userId int64) error
	IsRegistered(ctx context.Context, courseId int64, userId int64) (bool, error)
	DeleteCourse(ctx context.Context, courseId int64) error
	Stats(ctx context.Context) ([]entity.CourseStat, error)
}

type CourseTxRunner interface {
	AddCourseTransaction(ctx context.Context, txFunc func(ctx context.Context, tx AddCourseTx) error) error
	EditCourseTransaction(ctx context.Context, txFunc func(ctx context.Context, tx EditCourseTx) error) error
	ReorderLessonsTransaction(ctx context.Context, txFunc func(ctx context.Context, tx ReorderLessonsTx) error) error
}

type AddCourseTx interface {
	AddCourse(ctx context.Context, req entity.AddCourseRequest) error
}

type EditCourseTx interface {
	GetCoursePreviewPicture(ctx context.Context, id int64) (string, error)
	EditCourse(ctx context.Context, req entity.EditCourseRequest) error
}

type ReorderLessonsTx interface {
	UpdateLessonNumber(ctx context.Context, lessonId int64, number int64) error
}

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

func (s Course) ReorderLessons(ctx context.Context, req domain.EditCourseLessonsOrderingRequest) error {
	err := s.txRunner.ReorderLessonsTransaction(ctx, func(ctx context.Context, tx ReorderLessonsTx) error {
		for i, lessonId := range req.OrderedLessonsIds {
			err := tx.UpdateLessonNumber(ctx, lessonId, int64(i+1))
			if err != nil {
				return errors.WithMessagef(err, "update lesson number on %d lesson", lessonId)
			}
		}
		return nil
	})
	if err != nil {
		return errors.WithMessage(err, "reorder lessons tx")
	}
	return nil
}

func (s Course) Stats(ctx context.Context) ([]byte, error) {
	stats, err := s.courseRepo.Stats(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "stats")
	}
	toExport := [][]string{}
	toExport = append(toExport, []string{
		"\uFEFFидентификатор курса",
		"название курса",
		"ФИО автора курса",
		"количество зарегистрированных студентов",
	})

	for _, stat := range stats {
		toExport = append(toExport, []string{
			fmt.Sprint(stat.Id),
			stat.Title,
			stat.AuthorFio,
			fmt.Sprint(stat.Count),
		})
	}

	var b bytes.Buffer
	writer := csv.NewWriter(&b)
	err = writer.WriteAll(toExport)
	if err != nil {
		return nil, errors.WithMessage(err, "write all")
	}
	writer.Flush()
	return b.Bytes(), nil
}
