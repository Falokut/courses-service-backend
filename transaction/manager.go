package transaction

import (
	"context"

	"courses-service/repository"
	"courses-service/service"

	"github.com/Falokut/go-kit/client/db"
)

type Manager struct {
	db db.Transactional
}

func NewManager(db db.Transactional) *Manager {
	return &Manager{
		db: db,
	}
}

type loginTransaction struct {
	repository.Auth
	repository.User
}

func (m Manager) LoginTransaction(ctx context.Context, txFunc func(ctx context.Context, tx service.LoginTx) error) error {
	return m.db.RunInTransaction(ctx, func(ctx context.Context, tx *db.Tx) error {
		authRepo := repository.NewAuth(tx)
		userRepo := repository.NewUser(tx)
		return txFunc(ctx, loginTransaction{
			authRepo,
			userRepo,
		})
	})
}

type courseTransaction struct {
	repository.Course
}

func (m Manager) AddCourseTransaction(ctx context.Context, txFunc func(ctx context.Context, tx service.AddCourseTx) error) error {
	return m.db.RunInTransaction(ctx, func(ctx context.Context, tx *db.Tx) error {
		courseRepo := repository.NewCourse(tx)
		return txFunc(ctx, courseTransaction{
			courseRepo,
		})
	})
}
func (m Manager) EditCourseTransaction(ctx context.Context, txFunc func(ctx context.Context, tx service.EditCourseTx) error) error {
	return m.db.RunInTransaction(ctx, func(ctx context.Context, tx *db.Tx) error {
		courseRepo := repository.NewCourse(tx)
		return txFunc(ctx, courseTransaction{
			courseRepo,
		})
	})
}

type attachmentTransaction struct {
	repository.Attachment
}

func (m Manager) CleanAttachmentsTransaction(ctx context.Context, txFunc func(ctx context.Context, tx service.AttachmentCleanerTx) error) error {
	return m.db.RunInTransaction(ctx, func(ctx context.Context, tx *db.Tx) error {
		attachmentRepo := repository.NewAttachment(tx)
		return txFunc(ctx, attachmentTransaction{
			attachmentRepo,
		})
	})
}

type lessonTransaction struct {
	repository.Lesson
}

func (m Manager) AttachFileTransaction(ctx context.Context, txFunc func(ctx context.Context, tx service.AttachFileTx) error) error {
	return m.db.RunInTransaction(ctx, func(ctx context.Context, tx *db.Tx) error {
		lessonRepo := repository.NewLesson(tx)
		return txFunc(ctx, lessonTransaction{
			lessonRepo,
		})
	})
}
func (m Manager) DeleteVideoTransaction(ctx context.Context, txFunc func(ctx context.Context, tx service.DeleteVideoTx) error) error {
	return m.db.RunInTransaction(ctx, func(ctx context.Context, tx *db.Tx) error {
		lessonRepo := repository.NewLesson(tx)
		return txFunc(ctx, lessonTransaction{
			lessonRepo,
		})
	})
}

func (m Manager) AddVideoTransaction(ctx context.Context, txFunc func(ctx context.Context, tx service.AddVideoTx) error) error {
	return m.db.RunInTransaction(ctx, func(ctx context.Context, tx *db.Tx) error {
		lessonRepo := repository.NewLesson(tx)
		return txFunc(ctx, lessonTransaction{
			lessonRepo,
		})
	})
}

func (m Manager) ReorderLessonsTransaction(ctx context.Context, txFunc func(ctx context.Context, tx service.ReorderLessonsTx) error) error {
	return m.db.RunInTransaction(ctx, func(ctx context.Context, tx *db.Tx) error {
		lessonRepo := repository.NewLesson(tx)
		return txFunc(ctx, lessonTransaction{
			lessonRepo,
		})
	})
}
