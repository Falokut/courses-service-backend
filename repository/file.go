package repository

import (
	"context"
	"courses-service/entity"
	"fmt"
	"net/http"

	"github.com/Falokut/go-kit/http/client"
	"github.com/pkg/errors"
)

const (
	fileStorageEndpoint = "file"
)

type File struct {
	cli     *client.Client
	baseUrl string
}

func NewFile(cli *client.Client, baseUrl string) File {
	return File{
		cli:     cli,
		baseUrl: baseUrl,
	}
}

func (r File) GetFileUrl(category string, filename string) string {
	return fmt.Sprintf("%s/%s/%s/%s", r.baseUrl, fileStorageEndpoint, category, filename)
}

func (r File) UploadFile(ctx context.Context, req entity.UploadFileReq) error {
	url := r.GetFileUrl(req.Category, req.Filename)
	_, err := r.cli.Post(url).
		RequestBody(req.FileBody).
		StatusCodeToError().
		Do(ctx)
	if err != nil {
		return errors.WithMessagef(err, "call storage service %s", url)
	}
	return nil
}

func (r File) DeleteFile(ctx context.Context, url string) error {
	resp, err := r.cli.Delete(url).
		Do(ctx)
	if err != nil {
		return errors.WithMessagef(err, "call storage service %s", url)
	}

	switch {
	case resp.StatusCode() == http.StatusNotFound:
		return entity.ErrFileNotFound
	case !resp.IsSuccess():
		return errors.Errorf("call storage service %s, unexpected status: %d", url, resp.StatusCode())
	}
	return nil
}
