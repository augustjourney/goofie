package storage

import (
	"api/pkg/logger"
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type AWS struct {
	cfg    Config
	client *s3.S3
}

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

func (s *AWS) Upload(ctx context.Context, buffer io.ReadSeeker, bucket string, fileName string, mime string, expiry string) (string, error) {

	_, err := s.client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
		Body:   buffer,
	})

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
