package repository

import (
	"context"
	"courses-service/entity"

	"github.com/Falokut/go-kit/client/db"
	"github.com/pkg/errors"
)

type Lesson struct {
	db db.DB
}

func NewLesson(db db.DB) Lesson {
	return Lesson{
		db: db,
	}
}

func (r Lesson) CreateLesson(ctx context.Context, lesson entity.CreateLessonRequest) error {
	const query = `INSERT INTO course_lessons (course_id, lesson_number, created_at, title) VALUES($1, $2, $3, $4);`
	_, err := r.db.Exec(ctx, query, lesson.CourseId, lesson.LessonNumber, lesson.CreatedAt, lesson.Title)
	if err != nil {
		return errors.WithMessagef(err, "exec query: %s")
	}
	return nil
}

func (r Lesson) EditTitle(ctx context.Context, id int64, newTitle string) error {
	const query = "UPDATE course_lessons SET title=$1 WHERE id=$2;"
	_, err := r.db.Exec(ctx, query, newTitle, id)
	if err != nil {
		return errors.WithMessagef(err, "exec query: %s")
	}
	return nil
}

func (r Lesson) EditContent(ctx context.Context, id int64, content string) error {
	const query = "UPDATE course_lessons SET lesson_content=$1 WHERE id=$2;"
	_, err := r.db.Exec(ctx, query, content, id)
	if err != nil {
		return errors.WithMessagef(err, "exec query: %s")
	}
	return nil
}

func (r Lesson) UnattachFile(ctx context.Context, attachmentId int64) error {
	const query = "UPDATE lesson_attachments SET lesson_id=0 WHERE id=$1;"
	_, err := r.db.Exec(ctx, query, attachmentId)
	if err != nil {
		return errors.WithMessagef(err, "exec query: %s")
	}
	return nil
}

func (r Lesson) AttachFile(ctx context.Context, attachment entity.LessonAttachment) error {
	const query = "INSERT INTO lesson_attachments (lesson_id, attachment_type, pretty_name, url) VALUES($1, $2, $3, $4);"
	_, err := r.db.Exec(ctx, query, attachment.LessonId, attachment.Type, attachment.PrettyName, attachment.Url)
	if err != nil {
		return errors.WithMessagef(err, "exec query: %s")
	}
	return nil
}

func (r Lesson) DeleteVideo(ctx context.Context, id int64) (string, error) {
	const query = `
	WITH u AS (
    	SELECT video_url FROM course_lessons WHERE id = $1
	)
	UPDATE course_lessons SET video_url='' WHERE id=$1
	RETURNING (SELECT video_url FROM u);`

	url := ""
	err := r.db.SelectRow(ctx, &url, query, id)
	if err != nil {
		return "", errors.WithMessagef(err, "exec query: %s")
	}
	return url, nil
}

func (r Lesson) AddVideo(ctx context.Context, id int64, url string) error {
	const query = "UPDATE course_lessons SET video_url=$1 WHERE id=$2;"
	_, err := r.db.Exec(ctx, query, url, id)
	if err != nil {
		return errors.WithMessagef(err, "exec query: %s", query)
	}
	return nil
}

func (r Lesson) UpdateLessonNumber(ctx context.Context, lessonId int64, number int64) error {
	const query = "UPDATE course_lessons SET lesson_number=$1 WHERE id=$2;"
	_, err := r.db.Exec(ctx, query, number, lessonId)
	if err != nil {
		return errors.WithMessagef(err, "exec query: %s", query)
	}
	return nil
}
func (r Lesson) CheckLessonOwnership(ctx context.Context, userId int64, lessonId int64) (bool, error) {
	const query = `SELECT EXISTS (
		SELECT 1 
		FROM courses c 
		JOIN course_lessons cl ON c.id=cl.course_id 
		WHERE author_id=$1 AND cl.id=$2
		);`
	isOwner := false
	err := r.db.SelectRow(ctx, &isOwner, query, userId, lessonId)
	if err != nil {
		return false, errors.WithMessagef(err, "exec query: %s", query)
	}
	return isOwner, nil
}
