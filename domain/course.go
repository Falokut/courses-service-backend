package domain

import "time"

type GetCourseRequest struct {
	CourseId int64 `validate:"required"`
}

type CoursePreview struct {
	Id                int
	AuthorFio         string
	Title             string
	PreviewPictureUrl string
}

type Course struct {
	Id        int
	AuthorFio string
	Title     string
	Lessons   []Lesson
}

type Lesson struct {
	LessonNumber  int64
	Title         string
	CreatedAt     time.Time
	LessonContent string
	VideoUrl      string
}

type CourseIdRequest struct {
	CourseId int64 `validate:"required"`
}

type IsRegisteredResponse struct {
	IsRegistered bool
}
