package repository

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type FileRepository interface {
	Save(file multipart.File, fileName string) (string, error)
}

type fileRepository struct {
	path string
}

func (f *fileRepository) Save(file multipart.File, fileName string) (string, error) {
	fileLocation := filepath.Join(f.path, fileName)
	out, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}
	return fileLocation, nil
}

func NewFileRepository(path string) FileRepository {
	return &fileRepository{path: path}
}
