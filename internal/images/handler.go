package images

import (
	"api/pkg/config"
	"api/pkg/errs"
	"api/pkg/handler"
	"api/pkg/tracer"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	config  *config.Config
	service *Service
}

// Create handles http-request for creating and uploading an image
func (h *Handler) Create(c *fiber.Ctx) error {
	ctx := tracer.NewContext(c, "upload")

	var resp handler.Response

	userId, ok := handler.GetUserIDFromFiberContext(c)
	if !ok {
		return resp.WithStatus(fiber.StatusUnauthorized).Do(c)
	}

	file, err := handler.GetMultipartFormFile(ctx, c, "image")
	if err != nil {
		if errors.Is(err, errs.ErrNoMultipartFormData) {
			return resp.WithError(errs.ErrNoImageForUploading).Do(c)
		}
		return resp.WithError(err).Do(c)
	}

	err = h.service.Create(ctx, file, userId)
	if err != nil {
		return resp.WithError(err).Do(c)
	}

	return resp.Do(c)
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		config:  config.Get(),
		service: service,
	}
}
