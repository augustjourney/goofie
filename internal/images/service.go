package images

import (
	"api/internal/storage"
	"api/pkg/config"
	"api/pkg/files"
	"api/pkg/logger"
	"bytes"
	"context"
	"fmt"
	"image"
	"mime/multipart"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	config *config.Config
	repo   *Repo
	s3     storage.S3
}

func (s *Service) Create(ctx context.Context, file *multipart.FileHeader, authorID int) error {
	buff, err := openMultipart(ctx, file)
	if err != nil {
		return err
	}

	src, err := decode(ctx, buff)
	if err != nil {
		return err
	}

	// build image model
	img := Image{
		Ext:      files.GetExtension(ctx, file.Filename),
		Mime:     file.Header.Get("Content-Type"),
		Size:     file.Size,
		Name:     file.Filename,
		Slug:     uuid.New().String(),
		AuthorID: uint(authorID),
	}
	img.WithDefaults(s.config).WithMetadata(getMetadata(src))

	// upload original image
	filepath, err := s.s3.Upload(ctx, bytes.NewReader(buff.Bytes()), img.Bucket, img.GetFilename(), img.Mime, "")
	if err != nil {
		logger.Error(logger.Record{
			Error:   err,
			Context: ctx,
		})
		return err
	}
	fmt.Println("Done original ", filepath)

	// create image in db
	err = s.repo.Create(ctx, &img)
	if err != nil {
		// TODO: delete image from bucket
		// Because it was uploaded but not saved info to db
		return err
	}

	// pre-process image
	// create webp, jpeg, avif
	var rules []ResizeRule = []ResizeRule{
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
		go func(src *image.Image, r ResizeRule) {
			resized, err := resize(ctx, src, rule)
			if err != nil {
				return
			}
			expiryTime := time.Now().Add(time.Minute * 5)
			expiry := strconv.FormatInt(expiryTime.Unix(), 10)
			filename := getHashFilename(ctx, img.Slug, r)
			filepath, err := s.s3.Upload(ctx, bytes.NewReader(resized.Bytes()), img.Bucket, filename, "image/"+rule.Format, expiry)
			if err != nil {
				logger.Error(logger.Record{
					Error:   err,
					Context: ctx,
				})
				return
			}
			fmt.Println("Done ", filepath)
		}(src, rule)
	}

	return nil
}

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
