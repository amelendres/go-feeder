package cloud

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	"google.golang.org/api/drive/v3"
)

type GDFileProvider struct {
	drive *drive.Service
	file  *os.File
}

func NewGDFileProvider(ds *drive.Service) *GDFileProvider {

	return &GDFileProvider{ds, nil}
}

func (fp *GDFileProvider) File(fileId string) (*os.File, error) {

	file, err := fp.download(fileId)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (fp *GDFileProvider) download(fileId string) (*os.File, error) {

	//drive
	fgc := fp.drive.Files.Get(fileId)
	resp, err := fgc.Download()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	// defer resp.Body.Close()

	tmpFile, cleanTmpFile := createTempFile("import_tmp.docx")
	_, err = io.Copy(tmpFile, resp.Body)
	defer cleanTmpFile()

	if err != nil {
		log.Println("Error while downloading", fileId, "-", err)
		return nil, err
	}

	return tmpFile, nil
}

func createTempFile(fileName string) (*os.File, func()) {

	tmpFile, err := ioutil.TempFile("", fileName)

	if err != nil {
		log.Println("could not create temp file %v", err)
		return nil, nil
	}

	removeFile := func() {
		os.Remove(tmpFile.Name())
	}

	return tmpFile, removeFile
}
