package storage

import (
	"context"
	"io"
)

// Config stores data about connection to s3 storage.
type Config struct {
	Region      string
	Endpoint    string
	AccessKeyID string
	AccessKey   string
}

// S3 is the interface to the s3 storage interactions.
type S3 interface {
	Upload(ctx context.Context, buffer io.ReadSeeker, bucket string, fileName string, mime string, expiry string) (string, error)
}

// New creates and returns a new [S3] instance.
func New(provider string, cfg Config) S3 {
	if provider == "selectel" {
		return NewSelectel(cfg)
	}
	return NewAWS(cfg)
}

func buildFilePath(bucket, fileName string) string {
	return "/" + bucket + "/" + fileName
}
