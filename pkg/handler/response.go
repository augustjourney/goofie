package handler

import (
	"api/pkg/errs"
	"github.com/gofiber/fiber/v2"
)

// Response is a structure for response of all json http-requests body.
// StatusCode and Error not used in http response body.
type Response struct {
	OK         bool        `json:"ok"`
	Result     interface{} `json:"result"`
	Message    string      `json:"message,omitempty"`
	StatusCode int         `json:"-"`
	Error      error       `json:"-"`
}

// WithError sets error and status code to response.
func (r *Response) WithError(err error) *Response {
	r.Error = err

	// if status was already set with method `WithStatus`
	// do not override it
	if r.StatusCode != 0 {
		return r
	}

	status, ok := errs.Statuses[err]
	if ok {
		r.StatusCode = status
		return r
	}

	// if status not found in statuses map
	// use 500 status
	r.StatusCode = fiber.StatusInternalServerError

	return r
}

// WithStatus sets status code to response.
func (r *Response) WithStatus(status int) *Response {
	r.StatusCode = status
	return r
}

// WithData sets data to response result.
func (r *Response) WithData(data interface{}) *Response {
	r.Result = data
	return r
}

// WithMessage sets message to response.
func (r *Response) WithMessage(message string) *Response {
	r.Message = message
	return r
}

// Do does actual fiber response.
// If status code was not set, http.StatusOK will be used.
// If message was not set and error is not empty, error.Error() will be used as message.
func (r *Response) Do(fiberCtx *fiber.Ctx) error {
	if r.StatusCode == 0 {
		r.StatusCode = fiber.StatusOK
	}

	if r.StatusCode >= 400 || r.Error != nil {
		r.OK = false
	} else {
		r.OK = true
	}

	if r.Message == "" && r.Error != nil {
		r.Message = r.Error.Error()
	}

	if r.Message == "" && r.StatusCode == 500 {
		r.Message = "Internal Server Error"
	}

	return fiberCtx.Status(r.StatusCode).JSON(*r)
}
