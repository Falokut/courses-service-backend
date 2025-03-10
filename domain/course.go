package domain

type GetCourseRequest struct {
	CourseId int64 `validate:"required"`
}

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
	Lessons   []Lesson
	VideoUrl  string
}

type CourseIdRequest struct {
	CourseId int64 `validate:"required"`
}

type IsRegisteredResponse struct {
	IsRegistered bool
}

type AddCourseRequest struct {
	AuthorId       int64  `validate:"required"`
	Title          string `validate:"required"`
	PreviewPicture []byte `validate:"required"`
}

type AddCourseResponse struct {
	PreviewPictureUrl string
}

type EditCourseRequest struct {
	CourseId       int64  `validate:"required"`
	AuthorId       int64  `validate:"required"`
	Title          string `validate:"required"`
	PreviewPicture []byte `validate:"required"`
}

type EditCourseLessonsOrderingRequest struct {
	CourseId          int64   `validate:"required"`
	OrderedLessonsIds []int64 `validate:"required"`
}

type EditCourseResponse struct {
	NewPreviewPictureUrl string
}
