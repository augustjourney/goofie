package logger

import (
	"github.com/dusted-go/logging/prettylog"
	"log/slog"
	"runtime"
)

var (
	slogLogger *slog.Logger
)

var logLevels = map[slog.Level]string{
	slog.LevelDebug: "DEBUG",
	slog.LevelInfo:  "INFO",
	slog.LevelWarn:  "WARN",
	slog.LevelError: "ERROR",
}

func init() {
	prettyHandler := prettylog.NewHandler(&slog.HandlerOptions{
		Level:       slog.LevelDebug,
		AddSource:   false,
		ReplaceAttr: nil,
	})
	slogLogger = slog.New(prettyHandler)
}

func log(rec Record, level slog.Level) {
	rec.Type = logLevels[level]

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
func Info(rec Record) {
	log(rec, slog.LevelInfo)
}

// Error logs at Error level
func Error(rec Record) {
	log(rec, slog.LevelError)
}

// Debug logs at Debug level
func Debug(rec Record) {
	log(rec, slog.LevelDebug)
}

// Warn logs at Warn level
func Warn(rec Record) {
	log(rec, slog.LevelWarn)
}
