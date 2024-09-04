package storage

import (
	"api/pkg/logger"
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// AWS implements [S3] interface to interact with s3 aws storage.
type AWS struct {
	cfg    Config
	client *s3.S3
}

// NewAWS creates and returns a new [AWS] instance.
func NewAWS(cfg Config) *AWS {
	sess, _ := session.NewSession(&aws.Config{
		Region:           aws.String(cfg.Region),
		Endpoint:         aws.String(cfg.Endpoint),
		S3ForcePathStyle: aws.Bool(true),
		Credentials: credentials.NewStaticCredentials(
			cfg.AccessKeyID,
			cfg.AccessKey,
			"",
		),
	})

	return &AWS{
		cfg:    cfg,
		client: s3.New(sess),
	}
}

// Upload uploads a file to aws s3 storage.
func (s *AWS) Upload(ctx context.Context, buffer io.ReadSeeker, bucket string, fileName string, mime string, expiryTime *time.Duration) (string, error) {

	options := s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(fileName),
		Body:        buffer,
		ContentType: aws.String(mime),
	}

	if expiryTime != nil {
		expiry := time.Now().Add(*expiryTime)
		options.Expires = &expiry
	}

	_, err := s.client.PutObject(&options)

	if err == nil {
		return buildFilePath(bucket, fileName), nil
	}

	logger.Error(logger.Record{
		Message: "could not upload file to s3",
		Error:   err,
		Context: ctx,
		Data: map[string]interface{}{
			"Bucket":   bucket,
			"filename": fileName,
		},
	})

	return "", err
}
