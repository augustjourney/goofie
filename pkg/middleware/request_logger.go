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

func parseResponseBody(c *fiber.Ctx) interface{} {
	contentLength, _ := strconv.Atoi(c.Get("Content-Length"))

	if contentLength > 10000 {
		return "too big response body"
	}

	contentType := c.Get(fiber.HeaderContentType)

	if contentType == fiber.MIMETextPlain {
		return string(c.Response().Body())
	}

	if contentType != fiber.MIMEApplicationJSON {
		return fmt.Sprintf("unknown content type %s to process response body", contentType)
	}

	var responseBody interface{}

	err := json.Unmarshal(c.Response().Body(), &responseBody)
	if err != nil {
		return "unable to process response body"
	}

	return responseBody
}

func parseRequestBody(c *fiber.Ctx) interface{} {
	if c.Method() != http.MethodPost || c.Method() != http.MethodPut {
		return nil
	}

	var requestBody interface{}
	_ = c.BodyParser(&requestBody)
	return requestBody
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

	// get request payload
	// get request and response body
	data["request_body"] = parseRequestBody(c)
	data["response_body"] = parseResponseBody(c)

	ctx, ok := c.Locals("ctx").(context.Context)
	if !ok {
		ctx = context.Background()
	}

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
