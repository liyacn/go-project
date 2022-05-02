package storage

import (
	"context"
	"io"
)

type API interface {
	PutObject(ctx context.Context, path string, reader io.Reader) error
	DeleteObject(ctx context.Context, path string) error
	GetObject(ctx context.Context, path string) (io.ReadCloser, error)
	SignURL(ctx context.Context, path string, expiredInSec int64) (string, error)
}
