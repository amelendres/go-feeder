package fs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type StubReadsResource struct {
	file *os.File
}

func (s *StubReadsResource) Read(url string) (string, error) {
	//return "1\n\ntitle\n\npassage\n\ncontent\n\n2\n\ntitle\n\npassage\n\ncontent\n\n", nil
	return "one feed", nil
}

func TestDocFeeder(t *testing.T) {
	path := "Meditaciones 2019a.docx"

	// r := StubReadsResource{nil}
	fp := LocalFileProvider{}
	dp := DevotionalParser{}
	r := NewDocResource(&fp)
	df := NewDocFeeder(r, &dp)

	t.Run("it reads ten feeds", func(t *testing.T) {
		feeds, err := df.Feeds(path)

		assert.Empty(t, err)
		assert.Equal(t, 10, len(feeds))
	})
}
