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
