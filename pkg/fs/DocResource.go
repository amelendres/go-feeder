package fs

import (
	"fmt"
	"log"

	"code.sajari.com/docconv"
	feeder "github.com/amelendres/go-feeder/pkg"
)

type DocResource struct {
	fileProvider feeder.FileProvider
	content      string
}

func NewDocResource(fp feeder.FileProvider) *DocResource {

	return &DocResource{
		fileProvider: fp,
		content:      "",
	}
}

func (dr *DocResource) Read(url string) (string, error) {
	file, err := dr.fileProvider.File(url)
	if err != nil {
		log.Println(feeder.ErrOpeningFile, err)
		return "", fmt.Errorf("problem opening file  %v", err)
	}

	content, _, err := docconv.ConvertDocx(file)

	if err != nil {
		log.Println(feeder.ErrReadingFile, err)
		return "", fmt.Errorf("problem reading document %s, %v", file.Name(), err)
	}

	return content, nil
}
