package auth

import (
	"api/pkg/handler"
	"api/pkg/tracer"
	"github.com/gofiber/fiber/v2"
)

// Handler implements [IHandler] and stores methods for handling auth http-requests.
type Handler struct {
	service IService
}

// Signup handles http-request for user registration
func (h *Handler) Signup(c *fiber.Ctx) error {
	ctx := tracer.NewContext(c, "signup")

	var resp handler.Response
	var body SignupDTO

	err := c.BodyParser(&body)
	if err != nil {
		return resp.WithStatus(fiber.StatusBadRequest).WithError(err).Do(c)
	}

	err = body.Validate()
	if err != nil {
		return resp.WithStatus(fiber.StatusBadRequest).WithError(err).Do(c)
	}

	// do signup
	result, err := h.service.Signup(ctx, body)
	if err != nil {
		return resp.WithError(err).Do(c)
	}

	resp.Result = result

	if result.AlreadyExists {
		return resp.WithStatus(fiber.StatusConflict).WithMessage("You already have an account").Do(c)
	}

	return resp.WithMessage("Your account has been created successfully").Do(c)
}

// Login handles http-request for user login
func (h *Handler) Login(c *fiber.Ctx) error {
	ctx := tracer.NewContext(c, "login")
	var resp handler.Response
	var body LoginDTO

	err := c.BodyParser(&body)
	if err != nil {
		return resp.WithStatus(fiber.StatusBadRequest).WithError(err).Do(c)
	}

	err = body.Validate()
	if err != nil {
		return resp.WithStatus(fiber.StatusBadRequest).WithError(err).Do(c)
	}

	result, err := h.service.Login(ctx, body)
	if err != nil {
		return resp.WithError(err).Do(c)
	}

	resp.Result = result

	return resp.Do(c)
}

// NewHandler creates a new Auth Handler.
func NewHandler(service IService) *Handler {
	return &Handler{
		service,
	}
}
