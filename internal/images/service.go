package images

import (
	"api/internal/storage"
	"api/pkg/config"
	"api/pkg/logger"
	"bytes"
	"context"
	"mime/multipart"
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
	img := Image{}
	img.FromFileHeader(ctx, file).WithAuthor(authorID).WithDefaults(s.config).WithMetadata(getMetadata(src))

	// upload original image
	_, err = s.s3.Upload(ctx, bytes.NewReader(src), img.Bucket, img.GetFilename(), img.Mime, nil)
	if err != nil {
		logger.Error(logger.Record{
			Error:   err,
			Context: ctx,
			Message: "failed to upload original image",
			Data: map[string]interface{}{
				"image": img,
			},
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

	result.FromModel(img)

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
		go s.Resize(ctx, img, src, rule)
	}
}

func (s *Service) Resize(ctx context.Context, img Image, src []byte, rule ResizeRule) error {
	resized, err := resize(ctx, src, rule)
	if err != nil {
		return err
	}

	filename := getHashFilename(ctx, img.Slug, rule)
	mime := "image/" + rule.Format

	_, err = s.s3.Upload(ctx, bytes.NewReader(resized), img.Bucket, filename, mime, rule.ExpiryTime)
	if err != nil {
		logger.Error(logger.Record{
			Data: map[string]interface{}{
				"filename": filename,
				"image":    img,
				"rule":     rule,
			},
			Message: "failed to upload image",
			Error:   err,
			Context: ctx,
		})
		return err
	}

	return nil
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
