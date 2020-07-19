package fs

import (
	"fmt"

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
		return "", fmt.Errorf("problem reading file  %v", err)
	}

	content, _, err := docconv.ConvertDoc(file)

	//fmt.Println(content)

	if err != nil {
		return "", fmt.Errorf("problem reading document %s, %v", file.Name(), err)
	}

	return content, nil
}
