package fs

import (
	"code.sajari.com/docconv"
	"fmt"
	feed "github.com/amelendres/go-feeder/pkg"
)

type DocResource struct {
	fileProvider feed.FileProvider
	content      string
}

func NewDocResource(fp feed.FileProvider) feed.Reader {

	return &DocResource{
		fileProvider: fp,
		content:      "",
	}
}

func (dr *DocResource) Read(url string) (string, error) {
	file, err := dr.fileProvider.File(url)
	if err != nil {
		//log.Println(feed.ErrOpeningFile, err)
		return "", fmt.Errorf("Error opening file: %v ", err)
	}

	content, _, err := docconv.ConvertDoc(file)

	if err != nil {
		//log.Println(feed.ErrReadingFile, err)
		return "", fmt.Errorf("Error reading document: %s, %v ", file.Name(), err)
	}

	return content, nil
}
