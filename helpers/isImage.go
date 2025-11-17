package helpers

import (
	"io"
	"net/http"
	"path/filepath"
	"strings"
)

func IsImageExtension(fileName string) bool {
	ext := strings.ToLower(filepath.Ext(fileName))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
		return true
	default:
		return false
	}
}

func IsImageContent(file io.ReadSeeker) (bool, error) {
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return false, err
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return false, err
	}
	contentType := http.DetectContentType(buffer)
	return strings.HasPrefix(contentType, "image/"), nil
}
