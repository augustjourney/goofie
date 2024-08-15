package images

import (
	"api/pkg/config"
	"api/pkg/errs"
	"api/pkg/files"
	"api/pkg/handler"
	"api/pkg/tracer"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	config  *config.Config
	service *Service
}

// Upload handles http-request for creating and uploading an image
func (h *Handler) Upload(c *fiber.Ctx) error {
	ctx := tracer.NewContext(c, "upload")

	var resp handler.Response

	file, filePath, err := handler.SaveFile(ctx, c, "image")
	if err != nil {
		if errors.Is(err, errs.ErrNoMultipartFormData) {
			return resp.WithError(errs.ErrNoImageForUploading).Do(c)
		}
		return resp.WithError(err).Do(c)
	}

	defer files.Remove(ctx, filePath)

	var dto UploadImageDTO

	dto.Ext = files.GetExtension(ctx, file.Filename)
	dto.Name = file.Filename
	dto.Path = filePath
	dto.Size = file.Size
	dto.Mime = file.Header.Get("Content-Type")
	dto.AuthorID = 19

	err = h.service.Create(ctx, dto)
	if err != nil {
		return resp.WithError(err).Do(c)
	}

	return resp.Do(c)
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		config:  config.GetConfig(),
		service: service,
	}
}
