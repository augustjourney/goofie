package images

import (
	"api/pkg/logger"
	"bytes"
	"context"
	"image"
	"image/jpeg"

	"github.com/Kagami/go-avif"
	"github.com/chai2010/webp"
)

func toAvif(ctx context.Context, src *image.NRGBA, quality int) (*bytes.Buffer, error) {
	buff := new(bytes.Buffer)

	qualityPercent := quality * 100 / 100
	avifQuality := avif.MaxQuality - (avif.MaxQuality * qualityPercent / 100)

	err := avif.Encode(buff, src, &avif.Options{
		Quality: avifQuality,
	})

	if err == nil {
		return buff, nil
	}
	logger.Error(logger.Record{
		Context: ctx,
		Error:   err,
		Message: "could not encode image to avif",
	})
	return nil, err
}

func toWebp(ctx context.Context, src *image.NRGBA, quality int) (*bytes.Buffer, error) {
	buff := new(bytes.Buffer)
	err := webp.Encode(buff, src, &webp.Options{
		Quality: float32(quality),
	})
	if err == nil {
		return buff, nil
	}
	logger.Error(logger.Record{
		Context: ctx,
		Error:   err,
		Message: "could not encode image to webp",
	})
	return nil, err
}

func toJpeg(ctx context.Context, src *image.NRGBA, quality int) (*bytes.Buffer, error) {
	buff := new(bytes.Buffer)
	err := jpeg.Encode(buff, src, &jpeg.Options{
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
