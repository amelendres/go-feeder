package cloud

import (
	"io"
	"os"

	feed "github.com/amelendres/go-feeder/pkg"

	"google.golang.org/api/drive/v3"
)

type GDFileProvider struct {
	drive *drive.Service
	file  *os.File
}

func NewGDFileProvider(ds *drive.Service) feed.FileProvider {

	return &GDFileProvider{ds, nil}
}

func (fp *GDFileProvider) File(fileId string) (io.Reader, error) {

	file, err := fp.download(fileId)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (fp *GDFileProvider) download(fileId string) (io.Reader, error) {

	//drive
	fgc := fp.drive.Files.Get(fileId)
	resp, err := fgc.Download()
	if err != nil {
		return nil, err
	}
	// defer resp.Body.Close()

	return resp.Body, nil
}
