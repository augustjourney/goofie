package logger

import "log/slog"

var LogLevels map[slog.Level]string = map[slog.Level]string{
	slog.LevelDebug: "DEBUG",
	slog.LevelInfo:  "INFO",
	slog.LevelWarn:  "WARN",
	slog.LevelError: "ERROR",
	LevelRequest:    "REQUEST",
}

var LevelRequest slog.Level = slog.Level(-1)

func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.LevelKey {
		level := a.Value.Any().(slog.Level)
		logLevelString, ok := LogLevels[level]
		if ok {
			a.Value = slog.StringValue(logLevelString)
		}
	}
	return a
}
