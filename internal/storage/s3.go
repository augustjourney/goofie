package storage

import (
	"context"
	"io"
)

type Config struct {
	Region      string
	Endpoint    string
	AccessKeyID string
	AccessKey   string
}

type S3 interface {
	Upload(ctx context.Context, buffer io.ReadSeeker, bucket string, fileName string, mime string, expiry string) (string, error)
}

func New(provider string, cfg Config) S3 {
	if provider == "selectel" {
		return NewSelectel(cfg)
	}
	return NewAWS(cfg)
}

func buildFilePath(bucket, fileName string) string {
	return "/" + bucket + "/" + fileName
}
