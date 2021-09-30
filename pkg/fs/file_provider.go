package fs

import (
	"io"
	"os"

	feed "github.com/amelendres/go-feeder/pkg"
)

type FileProvider struct {
	file *os.File
}

func NewFileProvider() feed.FileProvider {
	return &FileProvider{}
}

func (fp *FileProvider) File(path string) (io.Reader, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}
	return file, nil
}

func (fp *FileProvider) Name() string {
	return "fs"
}
