package images

import (
	"api/internal/storage"
	"api/pkg/config"
	"api/pkg/files"
	"api/pkg/logger"
	"bytes"
	"context"
	"github.com/google/uuid"
	"mime/multipart"
	"strconv"
	"time"
)

// Service stores methods of images actions.
type Service struct {
	config *config.Config
	repo   *Repo
	s3     storage.S3
}

// Create creates a new image in the database and uploads it to s3 storage
func (s *Service) Create(ctx context.Context, file *multipart.FileHeader, authorID uint) (CreateResult, error) {
	var result CreateResult
	buff, err := openMultipart(ctx, file)
	if err != nil {
		return result, err
	}

	src := buff.Bytes()

	// build image model
	img := Image{
		Ext:      files.GetExtension(ctx, file.Filename),
		Mime:     file.Header.Get("Content-Type"),
		Size:     file.Size,
		Name:     file.Filename,
		Slug:     uuid.New().String(),
		AuthorID: authorID,
	}
	img.WithDefaults(s.config).WithMetadata(getMetadata(src))

	// upload original image
	_, err = s.s3.Upload(ctx, bytes.NewReader(src), img.Bucket, img.GetFilename(), img.Mime, "")
	if err != nil {
		logger.Error(logger.Record{
			Error:   err,
			Context: ctx,
		})
		return result, err
	}

	// create image in db
	err = s.repo.Create(ctx, &img)
	if err != nil {
		// TODO: delete image from bucket
		// Because it was uploaded but not saved info to db
		return result, err
	}

	s.ProcessUploadedImage(ctx, src, img)

	result.Slug = img.Slug
	result.Mime = img.Mime
	result.Ext = img.Ext
	result.Size = img.Size
	result.Width = img.Width
	result.Height = img.Height

	return result, nil
}

// ProcessUploadedImage converts uploaded image to webp, jpeg, avif and uploads new formats to s3 storage
func (s *Service) ProcessUploadedImage(ctx context.Context, src []byte, img Image) {
	var rules = []ResizeRule{
		{
			Quality: 80,
			Width:   img.Width,
			Height:  img.Height,
			Format:  "jpeg",
		},
		{
			Quality: 80,
			Width:   img.Width,
			Height:  img.Height,
			Format:  "webp",
		},
		{
			Quality: 80,
			Width:   img.Width,
			Height:  img.Height,
			Format:  "avif",
		},
	}

	for _, rule := range rules {
		go func(src []byte, r ResizeRule) {
			resized, err := resize(ctx, src, rule)
			if err != nil {
				return
			}

			expiryTime := time.Now().Add(time.Minute * 5)
			expiry := strconv.FormatInt(expiryTime.Unix(), 10)
			filename := getHashFilename(ctx, img.Slug, r)

			_, err = s.s3.Upload(ctx, bytes.NewReader(resized), img.Bucket, filename, "image/"+rule.Format, expiry)
			if err != nil {
				logger.Error(logger.Record{
					Error:   err,
					Context: ctx,
				})
				return
			}
		}(src, rule)
	}
}

// NewService creates and returns a new images [Service] instance.
func NewService(repo *Repo) *Service {
	cfg := config.Get()
	s3 := storage.New(cfg.S3Provider, storage.Config{
		Region:      cfg.S3Region,
		Endpoint:    cfg.S3Endpoint,
		AccessKeyID: cfg.S3AccessKeyId,
		AccessKey:   cfg.S3SecretAccessKey,
	})
	return &Service{
		config: cfg,
		repo:   repo,
		s3:     s3,
	}
}
