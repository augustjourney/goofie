package middleware

import (
	"api/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
	"time"
)

func parseResponseBody(c *fiber.Ctx) (interface{}, error) {
	contentLength, _ := strconv.Atoi(c.Get("Content-Length"))

	if contentLength > 10000 {
		return "too big response body", nil
	}

	contentType := c.Get(fiber.HeaderContentType)

	if contentType == fiber.MIMETextPlain {
		return string(c.Response().Body()), nil
	}

	if contentType != fiber.MIMEApplicationJSON {
		return fmt.Sprintf("unknown content type %s to process response body", contentType), nil
	}

	var responseBody interface{}

	err := json.Unmarshal(c.Response().Body(), &responseBody)
	if err != nil {
		return responseBody, err
	}

	return responseBody, nil
}

func parseRequestBody(c *fiber.Ctx) (interface{}, error) {
	if c.Method() != http.MethodPost || c.Method() != http.MethodPut {
		return nil, nil
	}

	var requestBody interface{}
	err := c.BodyParser(&requestBody)
	return requestBody, err
}

// RequestLogger logs incoming requests
func RequestLogger(c *fiber.Ctx) error {
	start := time.Now()
	result := c.Next()
	duration := time.Since(start).Milliseconds()

	// request data to log
	data := map[string]interface{}{
		"method":          c.Method(),
		"url":             c.OriginalURL(),
		"processing_time": duration,
		"status_code":     c.Response().StatusCode(),
		"query_params":    c.Queries(),
		"user_agent":      c.Get("user-agent"),
		"ip":              c.Get("X-Real-IP"),
		"path":            c.Path(),
	}

	ctx, ok := c.Locals("ctx").(context.Context)
	if !ok {
		ctx = context.Background()
	}

	// get request payload
	// get request and response body
	responseBody, err := parseResponseBody(c)
	if err != nil {
		logger.Error(ctx, "unable to parse response body", err)
	}

	requestBody, err := parseRequestBody(c)
	if err != nil {
		logger.Error(ctx, "unable to parse request body", err)
	}

	data["request_body"] = requestBody
	data["response_body"] = responseBody

	// log request data
	logData := logger.Record{
		Message: c.Path(),
		Data:    data,
		Type:    "REQUEST",
		Context: ctx,
	}
	logData.Log()

	return result
}
