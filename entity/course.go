package entity

import (
	"time"

	"github.com/Falokut/go-kit/json"
	"github.com/pkg/errors"
)

type CoursePreview struct {
	Id                int64
	AuthorId          int64
	AuthorFio         string
	Title             string
	PreviewPictureUrl string
}

type Course struct {
	Id        int64
	AuthorFio string
	Title     string
	Lessons   Lessons
}

type AddCourseRequest struct {
	AuthorId          int64
	Title             string
	PreviewPictureUrl string
}

type EditCourseRequest struct {
	Id                int64
	AuthorId          int64
	Title             string
	PreviewPictureUrl string
}

type Lesson struct {
	Id            int64
	LessonNumber  int64
	CourseId      int64
	Title         string
	CreatedAt     time.Time
	LessonContent string
	VideoUrl      string
	Attachments   []LessonAttachment
}

type LessonAttachment struct {
	Id         int64
	Type       string
	LessonId   int64
	PrettyName string
	Url        string
}

type Lessons []Lesson

func (o *Lessons) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.Errorf("failed to scan Lessons: %v", value)
	}
	return json.Unmarshal(bytes, o) //nolint:wrapcheck
}

type CourseStat struct {
	Id        int64
	Title     string
	AuthorFio string
	Count     int64
}
