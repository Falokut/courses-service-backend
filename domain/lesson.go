package domain

import "time"

type CreateLessonRequest struct {
	CourseId     int64  `validate:"required"`
	LessonNumber int64  `validate:"required"`
	Title        string `validate:"required"`
}

type LessonRequest struct {
	LessonId int64 `validate:"required"`
}

type EditLessonTitleRequest struct {
	LessonId int64  `validate:"required"`
	NewTitle string `validate:"required,min=5,max=100"`
}

type EditLessonContentRequest struct {
	LessonId   int64  `validate:"required"`
	NewContent string `validate:"required,min=10,max=1000"`
}

type AttachFileToLessonRequest struct {
	LessonId          int64  `validate:"required"`
	PrettyName        string `validate:"max=30"`
	AttachmentContent []byte `validate:"required"`
}

type UnattachFileRequest struct {
	LessonId     int64 `validate:"required"`
	AttachmentId int64 `validate:"required"`
}

type AddLessonVideoRequest struct {
	LessonId int64  `validate:"required"`
	Video    []byte `validate:"required"`
}

type Lesson struct {
	Id            int64
	LessonNumber  int64
	Title         string
	CreatedAt     time.Time
	LessonContent string
	VideoUrl      string
	Attachments   []LessonAttachment
}

type LessonAttachment struct {
	Id         int64
	LessonId   int64
	Type       string
	PrettyName string
	Url        string
}
