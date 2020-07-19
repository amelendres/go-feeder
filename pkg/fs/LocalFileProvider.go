package fs

import (
	"log"
	"os"
)

type LocalFileProvider struct {
	file *os.File
}

func (lfp *LocalFileProvider) File(path string) (*os.File, error) {
	file, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}
	return file, nil
}
