package fs

import (
	"testing"

	"github.com/amelendres/go-feeder/pkg/devom"
	"github.com/stretchr/testify/assert"
)

func TestDocFeeder(t *testing.T) {
	path := map[string]string{
		"feeds-10-0": "./_test_feeds-10-0.docx",
		"feeds-8-2":  "./_test_feeds-8-2.docx",
		"no-file":    "./_test_not-exist-file.docx",
		"2019a":      "./_test_2019a.docx",
	}

	fp := LocalFileProvider{}
	dp := devom.DevotionalParser{}
	r := NewDocResource(&fp)
	df := NewDocFeeder(r, &dp)

	t.Run("it reads 10 Feeds and 0 UnknownFeeds from Docx", func(t *testing.T) {
		feeds, unknownFeeds, err := df.Feeds(path["feeds-10-0"])

		assert.Empty(t, err)
		assert.Empty(t, unknownFeeds)
		assert.Equal(t, 10, len(feeds))
	})

	t.Run("it reads 8 Feeds and 2 UnknownFeeds from Docx", func(t *testing.T) {
		feeds, unknownFeeds, err := df.Feeds(path["feeds-8-2"])

		assert.Empty(t, err)
		assert.Equal(t, 8, len(feeds))
		assert.Equal(t, 2, len(unknownFeeds))
	})

	t.Run("it fails read feeds without resource file", func(t *testing.T) {
		feeds, unknownFeeds, err := df.Feeds(path["no-file"])

		assert.NotNil(t, err)
		assert.Equal(t, 0, len(feeds))
		assert.Equal(t, 0, len(unknownFeeds))
	})

	t.Run("it reads 2019a resource from Docx", func(t *testing.T) {
		feeds, unknownFeeds, err := df.Feeds(path["2019a"])

		assert.Nil(t, err)
		assert.Equal(t, 100, len(feeds))
		assert.Equal(t, 0, len(unknownFeeds))
	})
}
