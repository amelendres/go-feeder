package fs

import (
	feed "github.com/amelendres/go-feeder/pkg"
	"os"
)

type FileProvider struct {
	file *os.File
}

func NewFileProvider() feed.FileProvider  {
	return &FileProvider{}
}

func (fp *FileProvider) File(path string) (*os.File, error) {
	file, err := os.Open(path)

	if err != nil {
		// log.Fatal(err)
		return nil, err
	}
	return file, nil
}
