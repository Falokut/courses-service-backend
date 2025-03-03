package entity

import "github.com/pkg/errors"

const (
	CoursesCategory = "courses"
)

var (
	ErrFileNotFound = errors.New("file not found")
)

type UploadFileReq struct {
	Category string
	Filename string
	FileBody []byte
}
