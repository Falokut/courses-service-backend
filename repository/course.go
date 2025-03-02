package repository

import (
	"context"
	"courses-service/entity"

	"github.com/Falokut/go-kit/client/db"
	"github.com/pkg/errors"
)

type Course struct {
	db db.DB
}

func NewCourse(db db.DB) Course {
	return Course{
		db: db,
	}
}

func (r Course) GetCoursesPreview(ctx context.Context, limit int32, offset int32) ([]entity.CoursePreview, error) {
	const query = `
	SELECT c.id, u.fio AS author_fio, c.title, c.preview_picture_url
	FROM courses c
	JOIN users u ON c.author_id=u.id
	JOIN course_lessons l ON l.course_id=c.id
	WHERE c.id > 0
	GROUP BY c.id, u.fio, c.title, c.preview_picture_url
	ORDER BY c.id
	LIMIT $1 OFFSET $2;`

	coursePreview := make([]entity.CoursePreview, 0)
	err := r.db.Select(ctx, &coursePreview, query, limit, offset)
	if err != nil {
		return nil, errors.WithMessagef(err, "exec query: %s", query)
	}
	return coursePreview, nil
}

func (r Course) GetCourse(ctx context.Context, id int64) (entity.Course, error) {
	const query = `
	SELECT c.id, u.fio  AS author_fio, c.title,
		json_agg(
				json_build_object(
					'lessonNumber', l.lesson_number,
					'courseId', l.course_id,
					'title', l.title,
					'createdAt', l.created_at,
					'lessonContent', l.lesson_content,
					'videoUrl', l.video_url
				)
				) AS lessons
	FROM courses c
	JOIN users u ON c.author_id=u.id
	JOIN course_lessons l ON l.course_id=c.id
	WHERE c.id=$1
	GROUP BY c.id, u.fio, c.title, c.preview_picture_url;`

	var course entity.Course
	err := r.db.SelectRow(ctx, &course, query, id)
	if err != nil {
		return entity.Course{}, errors.WithMessagef(err, "exec query: %s", query)
	}
	return course, nil
}
