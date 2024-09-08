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

	logger.Error(ctx, "file does not contain extension", nil, "fileName", fileName, "ext", ext)
	return ""
}
