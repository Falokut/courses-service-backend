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

func (r Course) Register(ctx context.Context, courseId int64, userId int64) error {
	const query = `INSERT INTO courses_registration(course_id, user_id) VALUES($1,$2)
	ON CONFLICT DO NOTHING;`
	_, err := r.db.Exec(ctx, query, courseId, userId)
	if err != nil {
		return errors.WithMessagef(err, "exec query: %s", query)
	}
	return nil
}

func (r Course) IsRegistered(ctx context.Context, courseId int64, userId int64) (bool, error) {
	const query = `SELECT EXISTS (SELECT 1 FROM courses_registration WHERE course_id=$1 AND user_id=$2);`

	isRegistered := false
	err := r.db.SelectRow(ctx, &isRegistered, query, courseId, userId)
	if err != nil {
		return false, errors.WithMessagef(err, "exec query: %s", query)
	}
	return isRegistered, nil
}

func (r Course) GetUserCourses(ctx context.Context, userId int64) ([]entity.CoursePreview, error) {
	const query = `
	SELECT c.id, u.fio AS author_fio, c.title, c.preview_picture_url
	FROM courses c
	JOIN users u ON c.author_id=u.id
	JOIN course_lessons l ON l.course_id=c.id
	JOIN courses_registration cr ON c.id=cr.course_id
	WHERE cr.user_id=$1
	GROUP BY c.id, u.fio, c.title, c.preview_picture_url, cr.user_id
	ORDER BY c.id;`

	coursePreview := make([]entity.CoursePreview, 0)
	err := r.db.Select(ctx, &coursePreview, query, userId)
	if err != nil {
		return nil, errors.WithMessagef(err, "exec query: %s", query)
	}
	return coursePreview, nil
}
