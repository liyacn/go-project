package storage

import (
	"context"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type tencent struct {
	client *cos.Client
}

type CosConfig struct {
	BucketURL  string
	ServiceURL string
	SecretID   string
	SecretKey  string
}

func NewCOS(cfg *CosConfig) API {
	bu, _ := url.Parse(cfg.BucketURL)
	su, _ := url.Parse(cfg.ServiceURL)
	baseURL := &cos.BaseURL{
		BucketURL:  bu,
		ServiceURL: su,
	}
	client := cos.NewClient(baseURL, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cfg.SecretID,
			SecretKey: cfg.SecretKey,
		},
	})
	return &tencent{client}
}

func (s *tencent) PutObject(ctx context.Context, path string, reader io.Reader) error {
	_, err := s.client.Object.Put(ctx, strings.TrimLeft(path, "/"), reader, nil)
	return err
}

func (s *tencent) GetSignURL(ctx context.Context, path string, expiredInSec int64) (string, error) {
	auth := s.client.GetCredential()
	u, err := s.client.Object.GetPresignedURL(ctx, http.MethodGet, strings.TrimLeft(path, "/"),
		auth.SecretID, auth.SecretKey, time.Duration(expiredInSec)*time.Second, nil)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
