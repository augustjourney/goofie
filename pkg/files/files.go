package files

import (
	"api/pkg/logger"
	"context"
	"path/filepath"
)

// GetExtension extracts and returns file extension.
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
