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

// Upload handles http-request for creating and uploading an image
func (h *Handler) Upload(c *fiber.Ctx) error {
	ctx := tracer.NewContext(c, "upload")

	var resp handler.Response

	file, err := handler.GetMultipartFormFile(ctx, c, "image")
	if err != nil {
		if errors.Is(err, errs.ErrNoMultipartFormData) {
			return resp.WithError(errs.ErrNoImageForUploading).Do(c)
		}
		return resp.WithError(err).Do(c)
	}

	// TODO: get author id from token
	var AuthorID int = 19

	err = h.service.Create(ctx, file, AuthorID)
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
