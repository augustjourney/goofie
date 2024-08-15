package files

import (
	"api/pkg/logger"
	"context"
	"os"
	"path/filepath"
)

func Remove(ctx context.Context, path string) {
	err := os.Remove(path)
	if err == nil {
		return
	}

	logger.Error(logger.Record{
		Error:   err,
		Message: "could not removed local file",
		Context: ctx,
	})
}

func GetExtension(ctx context.Context, fileName string) string {
	ext := filepath.Ext(fileName)
	if ext != "" {
		return ext
	}

	logger.Error(logger.Record{
		Message: "file does not contain extension",
		Data: map[string]interface{}{
			"fileName": fileName,
			"ext":      ext,
		},
		Context: ctx,
	})
	return ""
}
