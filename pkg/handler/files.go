package handler

import (
	"api/pkg/config"
	"api/pkg/errs"
	"api/pkg/logger"
	"context"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/gofiber/fiber/v2"
)

func BuildTempFilePath(ctx context.Context, fileName string) string {
	wd, err := os.Getwd()
	if err != nil {
		logger.Error(logger.Record{
			Error:   err,
			Context: ctx,
			Message: "could not get wd",
		})
		return ""
	}

	cfg := config.GetConfig()

	return fmt.Sprintf("%s/%s/%s", wd, cfg.TempFolder, fileName)
}

func GetMultipartFormFile(c *fiber.Ctx, key string) (*multipart.FileHeader, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return nil, err
	}

	file := form.File[key][0]

	if file == nil {
		return nil, errs.ErrNoMultipartFormData
	}

	return file, nil
}

func SaveFile(ctx context.Context, c *fiber.Ctx, fileKey string) (*multipart.FileHeader, string, error) {
	file, err := GetMultipartFormFile(c, fileKey)
	if err != nil {
		logger.Error(logger.Record{
			Error:   err,
			Context: ctx,
			Message: "could not get multipart file",
			Data: map[string]interface{}{
				"fileKey": fileKey,
			},
		})
		return nil, "", err
	}

	filePath := BuildTempFilePath(ctx, file.Filename)

	err = c.SaveFile(file, filePath)
	if err == nil {
		return file, filePath, nil
	}

	logger.Error(logger.Record{
		Error:   err,
		Context: ctx,
		Message: "could not save file",
		Data: map[string]interface{}{
			"filePath": filePath,
			"fileName": file.Filename,
		},
	})

	return nil, "", errs.ErrInternal
}
