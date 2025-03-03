package entity

import (
	"time"

	"github.com/Falokut/go-kit/json"
	"github.com/pkg/errors"
)

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
	Lessons   Lessons
}

type Lesson struct {
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
