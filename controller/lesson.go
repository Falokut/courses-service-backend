package controller

import (
	"context"
	"courses-service/domain"
)

type LessonService interface {
	CreateLesson(ctx context.Context, req domain.CreateLessonRequest) error
	EditTitle(ctx context.Context, req domain.EditLessonTitleRequest) error
	EditLessonContent(ctx context.Context, req domain.EditLessonContentRequest) error
	AttachFileToLesson(ctx context.Context, req domain.AttachFileToLessonRequest) (string, error)
	UnattachFileFromLesson(ctx context.Context, req domain.UnattachFileRequest) error
	AddLessonVideo(ctx context.Context, req domain.AddLessonVideoRequest) (string, error)
	DeleteLessonVideo(ctx context.Context, id int64) error
}

type Lesson struct {
	service LessonService
}

func NewLesson(service LessonService) Lesson {
	return Lesson{
		service: service,
	}
}

// CreateLesson
//
//	@Tags		lesson
//	@Summary	Создать урок
//	@Accept		json
//	@Produce	json
//
//	@Param		body	body		domain.CreateLessonRequest	true	"тело запроса"
//
//	@Success	200		{object}	domain.EditCourseResponse
//	@Failure	400		{object}	apierrors.Error
//	@Failure	401		{object}	apierrors.Error
//	@Failure	403		{object}	apierrors.Error
//	@Failure	404		{object}	apierrors.Error
//	@Failure	500		{object}	apierrors.Error
//
//	@Security	Bearer
//
//	@Router		/lessons [POST]
func (c Lesson) CreateLesson(ctx context.Context, req domain.CreateLessonRequest) error {
	return c.service.CreateLesson(ctx, req)
}

// EditTitle
//
//	@Tags		lesson
//	@Summary	Редактировать заголовок урока
//	@Accept		json
//	@Produce	json
//
//	@Param		body	body		domain.EditLessonTitleRequest	true	"тело запроса"
//
//	@Success	200		{object}	any
//	@Failure	400		{object}	apierrors.Error
//	@Failure	401		{object}	apierrors.Error
//	@Failure	403		{object}	apierrors.Error
//	@Failure	404		{object}	apierrors.Error
//	@Failure	500		{object}	apierrors.Error
//
//	@Router		/lessons/edit_title [POST]
func (c Lesson) EditTitle(ctx context.Context, req domain.EditLessonTitleRequest) error {
	return c.service.EditTitle(ctx, req)
}

// EditLessonContent
//
//	@Tags		lesson
//	@Summary	Редактировать содержание урока
//	@Accept		json
//	@Produce	json
//
//	@Param		body	body		domain.EditLessonTitleRequest	true	"тело запроса"
//
//	@Success	200		{object}	any
//	@Failure	400		{object}	apierrors.Error
//	@Failure	401		{object}	apierrors.Error
//	@Failure	403		{object}	apierrors.Error
//	@Failure	404		{object}	apierrors.Error
//	@Failure	500		{object}	apierrors.Error
//
//	@Router		/lessons/edit_content [POST]
func (c Lesson) EditLessonContent(ctx context.Context, req domain.EditLessonContentRequest) error {
	return c.service.EditLessonContent(ctx, req)
}

// AttachFileToLesson
//
//	@Tags		lesson
//	@Summary	Прикрепить файл к уроку
//	@Accept		json
//	@Produce	json
//
//	@Param		body	body		domain.AttachFileToLessonRequest	true	"тело запроса"
//
//	@Success	200		{object}	any
//	@Failure	400		{object}	apierrors.Error
//	@Failure	401		{object}	apierrors.Error
//	@Failure	403		{object}	apierrors.Error
//	@Failure	404		{object}	apierrors.Error
//	@Failure	500		{object}	apierrors.Error
//
//	@Router		/lessons/attach_file [POST]
func (c Lesson) AttachFileToLesson(ctx context.Context, req domain.AttachFileToLessonRequest) (string, error) {
	return c.service.AttachFileToLesson(ctx, req)
}

// UnattachFileFromLesson
//
//	@Tags		lesson
//	@Summary	Открепить файл от урока
//	@Accept		json
//	@Produce	json
//
//	@Param		lessonId		query		int	true	"идентификатор"
//	@Param		attachmentId	query		int	true	"идентификатор"
//
//	@Success	200				{object}	any
//	@Failure	400				{object}	apierrors.Error
//	@Failure	401				{object}	apierrors.Error
//	@Failure	403				{object}	apierrors.Error
//	@Failure	404				{object}	apierrors.Error
//	@Failure	500				{object}	apierrors.Error
//
//	@Router		/lessons/unattach_file [DELETE]
func (c Lesson) UnattachFileFromLesson(ctx context.Context, req domain.UnattachFileRequest) error {
	return c.service.UnattachFileFromLesson(ctx, req)
}

// DeleteLessonVideo
//
//	@Tags		lesson
//	@Summary	Удалить видео из урока
//	@Accept		json
//	@Produce	json
//
//	@Param		body	body		domain.UnattachFileRequest	true	"тело запроса"
//
//	@Success	200		{object}	any
//	@Failure	400		{object}	apierrors.Error
//	@Failure	401		{object}	apierrors.Error
//	@Failure	403		{object}	apierrors.Error
//	@Failure	404		{object}	apierrors.Error
//	@Failure	500		{object}	apierrors.Error
//
//	@Router		/lessons/add_video [POST]
func (c Lesson) AddLessonVideo(ctx context.Context, req domain.AddLessonVideoRequest) (string, error) {
	return c.service.AddLessonVideo(ctx, req)
}

// DeleteLessonVideo
//
//	@Tags		lesson
//	@Summary	Удалить видео из урока
//	@Accept		json
//	@Produce	json
//
//	@Param		lessonId	query		int	true	"идентификатор"
//
//	@Success	200			{object}	any
//	@Failure	400			{object}	apierrors.Error
//	@Failure	401			{object}	apierrors.Error
//	@Failure	403			{object}	apierrors.Error
//	@Failure	404			{object}	apierrors.Error
//	@Failure	500			{object}	apierrors.Error
//
//	@Router		/lessons/delete_video [DELETE]
func (c Lesson) DeleteLessonVideo(ctx context.Context, req domain.LessonRequest) error {
	return c.service.DeleteLessonVideo(ctx, req.LessonId)
}
