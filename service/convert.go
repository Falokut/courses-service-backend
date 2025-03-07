package service

import (
	"cmp"
	"courses-service/domain"
	"courses-service/entity"
	"slices"
)

func entityCoursesPreviewToDomain(courses []entity.CoursePreview) []domain.CoursePreview {
	domainCourses := make([]domain.CoursePreview, 0, len(courses))
	for _, course := range courses {
		domainCourses = append(domainCourses, domain.CoursePreview{
			Id:                course.Id,
			AuthorId:          course.AuthorId,
			AuthorFio:         course.AuthorFio,
			Title:             course.Title,
			PreviewPictureUrl: course.PreviewPictureUrl,
		})
	}
	return domainCourses
}

func entityLessonsToDomain(lessons []entity.Lesson) []domain.Lesson {
	domainLessons := make([]domain.Lesson, 0, len(lessons))
	for _, lesson := range lessons {
		domainLessons = append(domainLessons, entityLessonToDomain(lesson))
	}
	slices.SortStableFunc(domainLessons, func(a, b domain.Lesson) int {
		return cmp.Compare(a.LessonNumber, b.LessonNumber) // sort asc
	})

	return domainLessons
}

func entityLessonToDomain(lesson entity.Lesson) domain.Lesson {
	attachments := entityLessonAttachmentsToDomain(lesson.Attachments)
	return domain.Lesson{
		Id:            lesson.Id,
		LessonNumber:  lesson.LessonNumber,
		Title:         lesson.Title,
		CreatedAt:     lesson.CreatedAt,
		LessonContent: lesson.LessonContent,
		Attachments:   attachments,
		VideoUrl:      lesson.VideoUrl,
	}
}

func entityLessonAttachmentsToDomain(attachments []entity.LessonAttachment) []domain.LessonAttachment {
	domainAttachments := make([]domain.LessonAttachment, 0, len(attachments))
	for _, attachment := range attachments {
		domainAttachments = append(domainAttachments, domain.LessonAttachment{
			Id:         attachment.Id,
			LessonId:   attachment.LessonId,
			Type:       attachment.Type,
			PrettyName: attachment.PrettyName,
			Url:        attachment.Url,
		})
	}
	return domainAttachments
}
