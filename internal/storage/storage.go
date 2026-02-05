package storage

import (
	"context"
	"io"
)

type Storage interface {
	Upload(ctx context.Context, file io.Reader, filename string, contentType string, size int64) (string, error)
	Download(ctx context.Context, path string) (io.ReadCloser, error)
	Delete(ctx context.Context, path string) error
}
