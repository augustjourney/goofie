package images

import (
	"api/pkg/config"
	"api/pkg/files"
	"api/pkg/logger"
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"image/png"

	"github.com/disintegration/imaging"
)

type Service struct {
	config  *config.Config
	storage *Repo
}

type Metadata struct {
	Width  int
	Height int
}

func (s *Service) GetMetadata(src *image.Image) Metadata {
	var metadata Metadata

	dimensions := (*src).Bounds().Size()

	metadata.Height = dimensions.Y
	metadata.Width = dimensions.X

	return metadata
}

// Open opens image from local path and returns type image.Image
func (s *Service) Open(ctx context.Context, path string) (*image.Image, error) {
	// read file from os
	src, err := imaging.Open(path)
	if err != nil {
		logger.Error(logger.Record{
			Error:   err,
			Context: ctx,
			Message: "failed to open image",
			Data: map[string]interface{}{
				"path": path,
			},
		})
		return nil, err
	}

	return &src, nil
}

// Encode converts image.Image to bytes.Buffer
func (s *Service) Encode(ctx context.Context, src *image.Image, mime string) (*bytes.Buffer, error) {
	buff := new(bytes.Buffer)
	var err error

	if mime == "image/png" {
		err = png.Encode(buff, *src)
	} else {
		err = jpeg.Encode(buff, *src, nil)
	}

	if err == nil {
		return buff, nil
	}

	logger.Error(logger.Record{
		Context: ctx,
		Error:   err,
		Message: "could not encode image",
		Data: map[string]interface{}{
			"mime": mime,
		},
	})

	return nil, err
}

func (s *Service) Upload(ctx context.Context, buff *bytes.Buffer, img *Image) error {

	cfg := files.UploadToS3Config{
		Region:      s.config.S3Region,
		Endpoint:    s.config.S3Endpoint,
		AccessKeyID: s.config.S3AccessKeyId,
		AccessKey:   s.config.S3SecretAccessKey,
	}

	// upload to originals
	_, err := files.UploadToS3(ctx, img.Bucket, img.GetFilename(), bytes.NewReader(buff.Bytes()), cfg)

	return err
}

func (s *Service) ToJpeg(ctx context.Context, src *image.Image, quality int) (*bytes.Buffer, error) {
	buff := new(bytes.Buffer)
	err := jpeg.Encode(buff, *src, &jpeg.Options{
		Quality: quality,
	})
	if err == nil {
		return buff, nil
	}
	logger.Error(logger.Record{
		Context: ctx,
		Error:   err,
		Message: "could not encode image to jpeg",
	})
	return nil, err
}

func (s *Service) Create(ctx context.Context, dto UploadImageDTO) error {
	// open image from path
	src, err := s.Open(ctx, dto.Path)
	if err != nil {
		return err
	}

	// build image model
	img := dto.ToModel().WithDefaults(s.config).WithMetadata(s.GetMetadata(src))

	// get image buffer
	buff, err := s.Encode(ctx, src, img.Mime)
	if err != nil {
		return err
	}

	// upload original image
	err = s.Upload(ctx, buff, img)
	if err != nil {
		return err
	}

	jpg, err := s.ToJpeg(ctx, src, 90)
	if err != nil {
		return err
	}

	img.Mime = "image/jpeg"
	img.Ext = "jpeg"

	err = s.Upload(ctx, jpg, img)
	if err != nil {
		return err
	}

	// create image in db
	err = s.storage.Create(ctx, img)

	return err
}

func NewService(storage *Repo) *Service {
	return &Service{
		config:  config.GetConfig(),
		storage: storage,
	}
}
