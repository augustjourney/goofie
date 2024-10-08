package logger

import (
	"context"
	"github.com/dusted-go/logging/prettylog"
	"log/slog"
	"runtime"
)

var (
	slogLogger *slog.Logger
)

func init() {
	prettyHandler := prettylog.NewHandler(&slog.HandlerOptions{
		Level:       slog.LevelDebug,
		AddSource:   false,
		ReplaceAttr: replaceAttr,
	})
	slogLogger = slog.New(prettyHandler)
}

func log(rec Record, level slog.Level) {
	rec.Type = LogLevels[level]

	if rec.Where == "" {
		pc, _, _, ok := runtime.Caller(2)
		details := runtime.FuncForPC(pc)

		if ok && details != nil {
			rec.Where = details.Name()
		}
	}

	rec.validate()
	rec.addValuesFromContext()

	args := make([]interface{}, 0)

	if len(rec.Data) > 0 {
		args = append(args, "data", rec.Data)
	}

	if rec.Type == "ERROR" {
		args = append(args, "error", rec.Error.Error())
	}

	args = append(args, "where", rec.Where)

	slogLogger.Log(rec.Context, level, rec.Message, args...)
}

// Info logs at Info level
func Info(ctx context.Context, message string, args ...interface{}) {
	log(Record{
		Context: ctx,
		Message: message,
		Data:    collectArgs(args...),
	}, slog.LevelInfo)
}

// Error logs at Error level
func Error(ctx context.Context, message string, err error, args ...interface{}) {
	log(Record{
		Context: ctx,
		Message: message,
		Error:   err,
		Data:    collectArgs(args...),
	}, slog.LevelError)
}

// Debug logs at Debug level
func Debug(ctx context.Context, message string, args ...interface{}) {
	log(Record{
		Context: ctx,
		Message: message,
		Data:    collectArgs(args...),
	}, slog.LevelDebug)
}

// Warn logs at Warn level
func Warn(ctx context.Context, message string, args ...interface{}) {
	log(Record{
		Context: ctx,
		Message: message,
		Data:    collectArgs(args...),
	}, slog.LevelWarn)
}
