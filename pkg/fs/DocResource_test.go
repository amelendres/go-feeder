package fs

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unidoc/unioffice/document"
)

type StubFileProvider struct {
	file *os.File
}

func (fp *StubFileProvider) File(path string) (*os.File, error) {
	docFile, cleanDocFile := createTempFile(path, "test devotional title")
	defer cleanDocFile()

	return docFile, nil
}

func createTempFile(fileName, initialData string) (*os.File, func()) {
	doc := document.New()
	doc.AddParagraph().AddRun().AddText(initialData)
	err := doc.SaveToFile(fileName)

	if err != nil {
		fmt.Printf("problem creating temporary file, %v", err)
	}

	removeFile := func() {
		os.Remove(fileName)
	}

	f, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("problem opening temporary file, %v", err)
	}

	return f, removeFile
}

func TestDocResource(t *testing.T) {
	text := "test devotional title"
	path := "test.doc"
	file, cleanDocFile := createTempFile(path, text)
	defer cleanDocFile()

	fp := StubFileProvider{file}
	dr := NewDocResource(&fp)

	t.Run("it reads a doc resource", func(t *testing.T) {

		fileContent, err := dr.Read(path)
		assert.Empty(t, err)
		assert.NotNil(t, fileContent)
		assert.Contains(t, fileContent, text)
	})
}
