package entity

import "time"

type CreateLessonRequest struct {
	CourseId     int64
	LessonNumber int64
	CreatedAt    time.Time
	Title        string
}
