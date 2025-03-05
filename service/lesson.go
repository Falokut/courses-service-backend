package service

import (
	"context"
	"courses-service/domain"
	"courses-service/entity"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type LessonRepo interface {
	CreateLesson(ctx context.Context, lesson entity.CreateLessonRequest) error
	EditTitle(ctx context.Context, id int64, newTitle string) error
	EditContent(ctx context.Context, id int64, content string) error
	UnattachFile(ctx context.Context, attachmentId int64) error
}

type LessonTxRunner interface {
	AttachFileTransaction(ctx context.Context, txFunc func(ctx context.Context, tx AttachFileTx) error) error
	DeleteVideoTransaction(ctx context.Context, txFunc func(ctx context.Context, tx DeleteVideoTx) error) error
	AddVideoTransaction(ctx context.Context, txFunc func(ctx context.Context, tx AddVideoTx) error) error
}

type AttachFileTx interface {
	AttachFile(ctx context.Context, attachment entity.LessonAttachment) error
}

type DeleteVideoTx interface {
	DeleteVideo(ctx context.Context, lessonId int64) (string, error)
}

type AddVideoTx interface {
	AddVideo(ctx context.Context, lessonId int64, url string) error
}

type Lesson struct {
	lessonRepo LessonRepo
	txRunner   LessonTxRunner
	fileRepo   FileRepo
}

func NewLesson(
	lessonRepo LessonRepo,
	txRunner LessonTxRunner,
	fileRepo FileRepo,
) Lesson {
	return Lesson{
		lessonRepo: lessonRepo,
		txRunner:   txRunner,
		fileRepo:   fileRepo,
	}
}

func (s Lesson) CreateLesson(ctx context.Context, req domain.CreateLessonRequest) error {
	err := s.lessonRepo.CreateLesson(ctx, entity.CreateLessonRequest{
		CourseId:     req.CourseId,
		LessonNumber: req.LessonNumber,
		CreatedAt:    time.Now().UTC(),
		Title:        req.Title,
	})
	if err != nil {
		return errors.WithMessage(err, "create lesson")
	}
	return nil
}

func (s Lesson) EditTitle(ctx context.Context, req domain.EditLessonTitleRequest) error {
	err := s.lessonRepo.EditTitle(ctx, req.LessonId, req.NewTitle)
	if err != nil {
		return errors.WithMessage(err, "edit lesson title")
	}
	return nil
}
func (s Lesson) EditLessonContent(ctx context.Context, req domain.EditLessonContentRequest) error {
	err := s.lessonRepo.EditContent(ctx, req.LessonId, req.NewContent)
	if err != nil {
		return errors.WithMessage(err, "edit lesson content")
	}
	return nil
}

func (s Lesson) AttachFileToLesson(ctx context.Context, req domain.AttachFileToLessonRequest) (string, error) {
	var url string
	err := s.txRunner.AttachFileTransaction(ctx, func(ctx context.Context, tx AttachFileTx) error {
		filename := uuid.NewString()
		url = s.fileRepo.GetFileUrl(entity.CoursesCategory, filename)
		err := tx.AttachFile(ctx, entity.LessonAttachment{
			Type:       "file",
			LessonId:   req.LessonId,
			PrettyName: req.PrettyName,
			Url:        url,
		})
		if err != nil {
			return errors.WithMessage(err, "attach file")
		}

		err = s.fileRepo.UploadFile(ctx, entity.UploadFileReq{
			Category: entity.CoursesCategory,
			Filename: filename,
			FileBody: req.AttachmentContent,
		})
		if err != nil {
			return errors.WithMessage(err, "upload file")
		}
		return nil
	})
	if err != nil {
		return "", errors.WithMessage(err, "attach file tx")
	}
	return url, nil
}

func (s Lesson) UnattachFileFromLesson(ctx context.Context, req domain.UnattachFileRequest) error {
	err := s.lessonRepo.UnattachFile(ctx, req.AttachmentId)
	if err != nil {
		return errors.WithMessage(err, "unattach file from lesson")
	}
	return nil
}

func (s Lesson) AddLessonVideo(ctx context.Context, req domain.AddLessonVideoRequest) (string, error) {
	var url string
	err := s.txRunner.AddVideoTransaction(ctx, func(ctx context.Context, tx AddVideoTx) error {
		filename := uuid.NewString()
		url = s.fileRepo.GetFileUrl(entity.CoursesCategory, filename)
		err := tx.AddVideo(ctx, req.LessonId, url)
		if err != nil {
			return errors.WithMessage(err, "add video")
		}

		err = s.fileRepo.UploadFile(ctx, entity.UploadFileReq{
			Category: entity.CoursesCategory,
			Filename: filename,
			FileBody: req.Video,
		})
		if err != nil {
			return errors.WithMessage(err, "upload video")
		}
		return nil
	})
	if err != nil {
		return "", errors.WithMessage(err, "add video tx")
	}
	return url, nil
}

func (s Lesson) DeleteLessonVideo(ctx context.Context, id int64) error {
	err := s.txRunner.DeleteVideoTransaction(ctx, func(ctx context.Context, tx DeleteVideoTx) error {
		url, err := tx.DeleteVideo(ctx, id)
		if err != nil {
			return errors.WithMessage(err, "delete video")
		}

		err = s.fileRepo.DeleteFile(ctx, url)
		if err != nil {
			return errors.WithMessage(err, "delete video")
		}
		return nil
	})
	if err != nil {
		return errors.WithMessage(err, "delete video tx")
	}
	return nil
}
