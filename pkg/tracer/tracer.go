package tracer

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	traceIDKey   = "trace_id"
	sessionIDKey = "session_id"
	spanNameKey  = "span_name"
)

func WithTraceID(ctx context.Context) context.Context {
	traceID, err := uuid.NewRandom()
	if err != nil {
		return context.WithValue(ctx, traceIDKey, "")
	}
	return context.WithValue(ctx, traceIDKey, traceID.String())
}

func GetTraceID(ctx context.Context) string {
	traceID, ok := ctx.Value(traceIDKey).(string)
	if !ok {
		return ""
	}
	return traceID
}

func WithSessionID(ctx context.Context, fiberCtx *fiber.Ctx) context.Context {
	// gets session id from request headers
	sessionID := fiberCtx.Get("session_id")
	return context.WithValue(ctx, sessionIDKey, sessionID)
}

func GetSessionID(ctx context.Context) string {
	sessionID, ok := ctx.Value(sessionIDKey).(string)
	if !ok {
		return ""
	}
	return sessionID
}

func WithSpanName(ctx context.Context, spanName string) context.Context {
	return context.WithValue(ctx, spanNameKey, spanName)
}

func GetSpanName(ctx context.Context) string {
	spanName, ok := ctx.Value(spanNameKey).(string)
	if !ok {
		return ""
	}
	return spanName
}

// NewContext creates a new context with generated trace id, session id and span name
func NewContext(fiberCtx *fiber.Ctx, spanName string) context.Context {
	ctx := WithTraceID(context.Background())
	ctx = WithSessionID(ctx, fiberCtx)
	ctx = WithSpanName(ctx, spanName)

	// put this context to fiber context
	// then we can get it like so: fiberCtx.Locals("ctx").(context.Context)
	// currently it is used in request middleware to extract values above and log them
	fiberCtx.Locals("ctx", ctx)

	// put trace_id to response headers
	fiberCtx.Set("trace_id", GetTraceID(ctx))

	return ctx
}
