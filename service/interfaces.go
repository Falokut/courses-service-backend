package service

import (
	"context"
	"courses-service/entity"
)

type UserRepo interface {
	Register(ctx context.Context, req entity.RegisterUser) error
	UpsertUser(ctx context.Context, req entity.UpsertUser) error
	GetUsers(ctx context.Context, limit int32, offset int32) ([]entity.User, error)
	DeleteUser(ctx context.Context, userId int32) error
	GetUserBySessionId(ctx context.Context, sessionId string) (*entity.User, error)
	UpdateUser(ctx context.Context, req entity.User) error
	GetUsersByRoleName(ctx context.Context, roleName string) ([]entity.User, error)
}

type AuthTxRunner interface {
	LoginTransaction(ctx context.Context, txFunc func(ctx context.Context, tx LoginTx) error) error
}

type LoginTx interface {
	GetUserByUsername(ctx context.Context, username string) (entity.User, error)
	InsertSession(ctx context.Context, session entity.Session) error
}

type AuthRepo interface {
	GetUserSession(ctx context.Context, sessionId string) (entity.UserSession, error)
	DeleteSession(ctx context.Context, sessionId string) error
	DeleteUserSession(ctx context.Context, sessionId string, userId int64) error
	GetUserSessions(ctx context.Context, userId int64) ([]entity.Session, error)
}

type RoleRepo interface {
	GetRoleId(ctx context.Context, roleName string) (int32, error)
	GetRoles(ctx context.Context) ([]entity.Role, error)
}

type CourseTxRunner interface {
	AddCourseTransaction(ctx context.Context, txFunc func(ctx context.Context, tx AddCourseTx) error) error
	EditCourseTransaction(ctx context.Context, txFunc func(ctx context.Context, tx EditCourseTx) error) error
}

type AddCourseTx interface {
	AddCourse(ctx context.Context, req entity.AddCourseRequest) error
}

type EditCourseTx interface {
	GetCoursePreviewPicture(ctx context.Context, id int64) (string, error)
	EditCourse(ctx context.Context, req entity.EditCourseRequest) error
}

type CourseRepo interface {
	GetCoursesPreview(ctx context.Context, limit int32, offset int32) ([]entity.CoursePreview, error)
	GetCourse(ctx context.Context, id int64) (entity.Course, error)
	GetUserCourses(ctx context.Context, userId int64) ([]entity.CoursePreview, error)
	GetCoursesByAuthorId(ctx context.Context, userId int64) ([]entity.CoursePreview, error)
	Register(ctx context.Context, courseId int64, userId int64) error
	IsRegistered(ctx context.Context, courseId int64, userId int64) (bool, error)
	DeleteCourse(ctx context.Context, courseId int64) error
}

type FileRepo interface {
	UploadFile(ctx context.Context, req entity.UploadFileReq) error
	DeleteFile(ctx context.Context, url string) error
	GetFileUrl(category string, filename string) string
}
