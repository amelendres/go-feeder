package cloud

import (
	"fmt"
	"io"
	"os"
	"regexp"

	feed "github.com/amelendres/go-feeder/pkg"

	"google.golang.org/api/drive/v3"
)

var ErrNotFoundGoogleDriveFileId = func(url string) error {
	return fmt.Errorf("Url <%s> does not have the file id", url)
}

type GDFileProvider struct {
	drive *drive.Service
	file  *os.File
}

func NewGDFileProvider(ds *drive.Service) feed.FileProvider {

	return &GDFileProvider{ds, nil}
}

func (fp *GDFileProvider) File(url string) (io.Reader, error) {

	fileId, err := fp.fileId(url)
	if err != nil {
		return nil, err
	}

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

func (fp *GDFileProvider) Name() string {
	return "gd"
}

func (fp *GDFileProvider) fileId(url string) (string, error) {
	re := regexp.MustCompile(`[0-9A-Za-z_-]{33}`)
	id := re.FindAllString(url, 1)
	if len(id) == 1 {
		return id[0], nil
	}
	return "", ErrNotFoundGoogleDriveFileId(url)
}
