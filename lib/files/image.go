package files

import (
	"github.com/blazee5/quizmaster-backend/lib/http_errors"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

func PrepareImage(fileHeader *multipart.FileHeader) (string, []byte, string, error) {
	file, err := fileHeader.Open()
	defer file.Close()

	if err != nil {
		return "", nil, "", err
	}

	bytes, err := io.ReadAll(file)

	if err != nil {
		return "", nil, "", err
	}

	contentType := http.DetectContentType(bytes)

	if !checkImageMime(contentType) {
		return "", nil, "", http_errors.ErrInvalidImage
	}

	id, err := uuid.NewUUID()

	if err != nil {
		return "", nil, "", err
	}

	fileName := id.String() + filepath.Ext(fileHeader.Filename)

	return contentType, bytes, fileName, nil
}

func checkImageMime(imageMime string) bool {
	var imageMimeTypes = map[string]struct{}{
		"image/gif":  {},
		"image/jpeg": {},
		"image/png":  {},
		"image/webp": {},
	}

	_, ok := imageMimeTypes[imageMime]
	return ok
}
