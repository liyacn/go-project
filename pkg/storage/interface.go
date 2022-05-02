package storage

import (
	"context"
	"io"
)

type API interface {
	PutObject(ctx context.Context, path string, reader io.Reader) error
	GetSignURL(ctx context.Context, path string, expiredInSec int64) (string, error)
}
