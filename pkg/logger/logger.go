package logger

import (
	"log/slog"
	"os"
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
	slogLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
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

func Info(rec Record) {
	log(rec, slog.LevelInfo)
}

func Error(rec Record) {
	log(rec, slog.LevelError)
}

func Debug(rec Record) {
	log(rec, slog.LevelDebug)
}

func Warn(rec Record) {
	log(rec, slog.LevelWarn)
}
