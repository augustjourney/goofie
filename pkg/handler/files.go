package handler

import (
	"api/pkg/errs"
	"api/pkg/logger"
	"context"
	"github.com/gofiber/fiber/v2"
	"mime/multipart"
)

// GetMultipartFormFile extracts uploading file from multipart form by given fileKey.
func GetMultipartFormFile(ctx context.Context, c *fiber.Ctx, fileKey string) (*multipart.FileHeader, error) {
	form, err := c.MultipartForm()
	if err != nil {
		logger.Error(ctx, "could not get multipart file", err, "fileKey", fileKey)
		return nil, err
	}

	if len(form.File[fileKey]) == 0 {
		return nil, errs.ErrNoMultipartFormData
	}

	file := form.File[fileKey][0]

	if file == nil {
		return nil, errs.ErrNoMultipartFormData
	}

	return file, nil
}
