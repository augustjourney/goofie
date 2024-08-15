package files

import (
	"context"
	"io"

	"api/pkg/errs"
	"api/pkg/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type UploadToS3Config struct {
	Region      string
	Endpoint    string
	AccessKeyID string
	AccessKey   string
}

func newS3Client(ctx context.Context, config UploadToS3Config) (*s3.S3, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(config.Region),
		Endpoint:         aws.String(config.Endpoint),
		S3ForcePathStyle: aws.Bool(true),
		Credentials: credentials.NewStaticCredentials(
			config.AccessKeyID,
			config.AccessKey,
			"",
		),
	})

	if err != nil {
		logger.Error(logger.Record{
			Message: "could not create session in aws sdk",
			Error:   err,
			Context: ctx,
		})
		return nil, err
	}

	return s3.New(sess), nil
}

func UploadToS3(ctx context.Context, bucket string, fileName string, fileBytes io.ReadSeeker, config UploadToS3Config) (*s3.PutObjectOutput, error) {
	client, err := newS3Client(ctx, config)
	if err != nil {
		return nil, errs.ErrInternal
	}

	result, err := client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
		Body:   fileBytes,
	})

	if err == nil {
		return result, nil
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

	return nil, err
}
