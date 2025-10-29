package helpers

import (
	"path/filepath"
	"strings"
)

func IsImage(fileName string) bool {
	ext := strings.ToLower(filepath.Ext(fileName))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
		return true
	default:
		return false
	}
}
