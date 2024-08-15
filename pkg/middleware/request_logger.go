package middleware

import (
	"api/pkg/logger"
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"time"
)

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
	var requestBody interface{}
	err := c.BodyParser(&requestBody)
	if err != nil {
		logger.Error(logger.Record{
			Message: "[middleware: RequestLogger] Could not parse request body",
			Error:   err,
			Data:    data,
		})
	}
	data["request_body"] = requestBody

	// get response body
	contentLength, _ := strconv.Atoi(c.Get("Content-Length"))
	var responseBody interface{}
	if contentLength < 10000 {
		err = json.Unmarshal(c.Response().Body(), &responseBody)
		if err != nil {
			logger.Error(logger.Record{
				Message: "[middleware: RequestLogger] Could not parse response body",
				Error:   err,
				Data:    data,
			})
		}
	} else {
		responseBody = []byte("Too big response body")
	}

	data["response_body"] = responseBody

	ctx, ok := c.Locals("ctx").(context.Context)
	if !ok {
		ctx = context.Background()
	}

	// log request data
	logger.Info(logger.Record{
		Message: c.Path(),
		Data:    data,
		Type:    "REQUEST",
		Context: ctx,
	})

	return result
}
