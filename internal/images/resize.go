package images

import (
	"api/pkg/consts"
	"api/pkg/errs"
	"api/pkg/logger"
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"github.com/h2non/bimg"
	"io"
	"mime/multipart"
	"time"
)

const (
	// MinQuality is the smallest number of [ResizeRule.Quality]
	MinQuality = 0
	// MaxQuality is the largest number of [ResizeRule.Quality]
	MaxQuality = 100
)

// ResizeRule stores data how image should be resized and converted.
type ResizeRule struct {
	Quality int    `json:"quality"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Format  string `json:"format"`
}

// Validate checks values of Quality, Height and Width.
func (r *ResizeRule) Validate() error {
	if r.Quality < MinQuality || r.Quality > MaxQuality {
		return errs.ErrWrongQualityValue
	}

	if r.Height < 0 {
		return errs.ErrWrongHeightValue
	}

	if r.Width < 0 {
		return errs.ErrWrongWidthValue
	}

	return nil
}

func getHashFilename(ctx context.Context, slug string, rule ResizeRule) string {
	name := fmt.Sprintf("%s_%d_%d_%d_%s", slug, rule.Width, rule.Height, rule.Quality, rule.Format)
	hash := md5.New()
	_, err := io.WriteString(hash, name)
	if err != nil {
		logger.Error(logger.Record{
			Error:   err,
			Context: ctx,
			Message: "failed to hash filename",
			Data: map[string]interface{}{
				"slug": slug,
				"rule": rule,
			},
		})
	}
	if rule.Format == "" {
		rule.Format = "jpeg"
	}
	return fmt.Sprintf("%x", hash.Sum(nil)) + "." + rule.Format
}

func resize(ctx context.Context, src []byte, rule ResizeRule) ([]byte, error) {
	startResizingAt := time.Now()
	options := bimg.Options{
		Width:   rule.Width,
		Height:  rule.Height,
		Quality: rule.Quality,
	}

	switch rule.Format {
	case "jpeg":
		options.Type = bimg.JPEG
	case "webp":
		options.Type = bimg.WEBP
	case "avif":
		options.Type = bimg.AVIF
	default:
		return nil, errs.ErrOutputFormatNotSupported
	}

	result, err := bimg.NewImage(src).Process(options)
	resizeProcessingTime := time.Since(startResizingAt).Milliseconds()
	ctx = context.WithValue(ctx, consts.ResizeProcessingTimeKey, resizeProcessingTime)
	if err != nil {
		logger.Error(logger.Record{
			Error:   err,
			Context: ctx,
		})
		return nil, err
	}

	logger.Info(logger.Record{
		Message: fmt.Sprintf("resize image success, format: %s, quality: %d, width: %d, height: %d", rule.Format, rule.Quality, rule.Width, rule.Height),
		Context: ctx,
	})

	return result, err

}

// Metadata stores data about uploaded image
type Metadata struct {
	Width  int
	Height int
}

func getMetadata(src []byte) Metadata {
	var metadata Metadata

	result, err := bimg.NewImage(src).Size()
	if err != nil {
		logger.Error(logger.Record{
			Error:   err,
			Message: "failed to get metadata",
		})
		return metadata
	}

	metadata.Width = result.Width
	metadata.Height = result.Height

	return metadata
}

func openMultipart(ctx context.Context, file *multipart.FileHeader) (*bytes.Buffer, error) {
	fl, err := file.Open()
	if err != nil {
		logger.Error(logger.Record{
			Error:   err,
			Context: ctx,
			Message: "failed to open image",
			Data: map[string]interface{}{
				"file": file,
			},
		})
		return nil, err
	}

	flRead, err := io.ReadAll(fl)
	if err != nil {
		logger.Error(logger.Record{
			Error:   err,
			Context: ctx,
			Message: "failed to read image",
			Data: map[string]interface{}{
				"file": file,
			},
		})
		return nil, err
	}

	return bytes.NewBuffer(flRead), nil
}

// TODO
func exportAsBase64(ctx context.Context) {}

// TODO
func saveToFile(ctx context.Context) {}
