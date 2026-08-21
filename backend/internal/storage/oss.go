package storage

import (
	"context"
	"fmt"
	"io"
	"strings"

	goliboss "codeup.aliyun.com/61e54b0e0bb300d827e1ae27/backend/golib/oss"
	"codezone/backend/internal/config"
	"github.com/spf13/viper"
)

type PublicStorage interface {
	Upload(ctx context.Context, key string, content io.Reader, contentType string) error
	URL(key string) string
}

type OSS struct {
	bucketName string
	urlPrefix  string
}

func NewOSS(cfg config.OSS) (*OSS, error) {
	if cfg.Endpoint == "" || cfg.AccessKeyID == "" || cfg.AccessKeySecret == "" || cfg.BucketName == "" || cfg.URL == "" {
		return nil, nil
	}
	viper.Set("oss.end_point", cfg.Endpoint)
	viper.Set("oss.access_key_ID", cfg.AccessKeyID)
	viper.Set("oss.access_key_secret", cfg.AccessKeySecret)
	if err := goliboss.Init(); err != nil {
		return nil, fmt.Errorf("initialize oss: %w", err)
	}
	return &OSS{bucketName: cfg.BucketName, urlPrefix: strings.TrimRight(cfg.URL, "/")}, nil
}

func (s *OSS) Upload(_ context.Context, key string, content io.Reader, contentType string) error {
	if _, err := goliboss.GetClient().Upload(s.bucketName, key, content, contentType); err != nil {
		return fmt.Errorf("upload avatar: %w", err)
	}
	return nil
}

func (s *OSS) URL(key string) string {
	return s.urlPrefix + "/" + strings.TrimLeft(key, "/")
}
