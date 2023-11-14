package http_utils

import (
	"io"
	"mime/multipart"
	"os"
)

func UploadFile(file *multipart.FileHeader, destinationPath string) error {
	if _, err := os.Stat("public"); os.IsNotExist(err) {
		_ = os.Mkdir("public", os.ModePerm)
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}
