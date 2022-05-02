package storage

import (
	"context"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"log"
	"strings"
)

type ali struct {
	bucket *oss.Bucket
}

type OssConfig struct {
	Endpoint   string
	KeyID      string
	KeySecret  string
	BucketName string
}

func NewOSS(cfg *OssConfig) API {
	client, err := oss.New(cfg.Endpoint, cfg.KeyID, cfg.KeySecret)
	if err != nil {
		log.Fatal(err)
	}
	bucket, err := client.Bucket(cfg.BucketName)
	if err != nil {
		log.Fatal(err)
	}
	return &ali{bucket}
}

func (s *ali) PutObject(ctx context.Context, path string, reader io.Reader) error {
	return s.bucket.PutObject(strings.TrimLeft(path, "/"), reader, oss.WithContext(ctx))
}

func (s *ali) DeleteObject(ctx context.Context, path string) error {
	return s.bucket.DeleteObject(strings.TrimLeft(path, "/"), oss.WithContext(ctx))
}

func (s *ali) GetObject(ctx context.Context, path string) (io.ReadCloser, error) {
	return s.bucket.GetObject(strings.TrimLeft(path, "/"), oss.WithContext(ctx))
}

func (s *ali) SignURL(ctx context.Context, path string, expiredInSec int64) (string, error) {
	return s.bucket.SignURL(strings.TrimLeft(path, "/"), oss.HTTPGet, expiredInSec, oss.WithContext(ctx))
}
