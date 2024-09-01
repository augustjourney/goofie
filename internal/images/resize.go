package images

import (
	"api/pkg/errs"
	"api/pkg/logger"
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"image"
	"io"
	"mime/multipart"

	"github.com/disintegration/imaging"
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

func resize(ctx context.Context, src *image.Image, rule ResizeRule) (*bytes.Buffer, error) {
	resizedImage := imaging.Resize(*src, rule.Width, rule.Height, imaging.Lanczos)

	var buff *bytes.Buffer
	var err error

	switch rule.Format {
	case "jpeg":
		buff, err = toJpeg(ctx, resizedImage, rule.Quality)
	case "webp":
		buff, err = toWebp(ctx, resizedImage, rule.Quality)
	case "avif":
		buff, err = toAvif(ctx, resizedImage, rule.Quality)
	default:
		return nil, errs.ErrOutputFormatNotSupported
	}

	return buff, err

}

// Metadata stores data about uploaded image
type Metadata struct {
	Width  int
	Height int
}

func getMetadata(src *image.Image) Metadata {
	var metadata Metadata

	dimensions := (*src).Bounds().Size()

	metadata.Height = dimensions.Y
	metadata.Width = dimensions.X

	return metadata
}

func decode(ctx context.Context, file *bytes.Buffer) (*image.Image, error) {
	src, err := imaging.Decode(bytes.NewReader(file.Bytes()))
	if err != nil {
		logger.Error(logger.Record{
			Error:   err,
			Context: ctx,
			Message: "failed to decode image",
			Data: map[string]interface{}{
				"file": file,
			},
		})
		return nil, err
	}
	return &src, nil
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
