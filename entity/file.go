package entity

const (
	CoursesCategory = "courses"
)

type UploadFileReq struct {
	Category string
	Filename string
	FileBody []byte
}
