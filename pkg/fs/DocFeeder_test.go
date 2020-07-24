package fs

import (
	"testing"

	"github.com/amelendres/go-feeder/pkg/devom"
	"github.com/stretchr/testify/assert"
)

func TestDocFeeder(t *testing.T) {
	path := []string{
		"./_test_feeds-10-0.docx",
		"./_test_feeds-8-2.docx",
		"./_test_not-exist-file.docx",
	}

	fp := LocalFileProvider{}
	dp := devom.DevotionalParser{}
	r := NewDocResource(&fp)
	df := NewDocFeeder(r, &dp)

	t.Run("it reads 10 Feeds and 0 UnknownFeeds", func(t *testing.T) {
		feeds, unknownFeeds, err := df.Feeds(path[0])

		assert.Empty(t, err)
		assert.Empty(t, unknownFeeds)
		assert.Equal(t, 10, len(feeds))
	})

	t.Run("it reads 8 Feeds and 2 UnknownFeeds", func(t *testing.T) {
		feeds, unknownFeeds, err := df.Feeds(path[1])

		assert.Empty(t, err)
		assert.Equal(t, 8, len(feeds))
		assert.Equal(t, 2, len(unknownFeeds))
	})

	t.Run("it fails read feeds without file resource", func(t *testing.T) {
		feeds, unknownFeeds, err := df.Feeds(path[2])

		assert.NotNil(t, err)
		assert.Equal(t, 0, len(feeds))
		assert.Equal(t, 0, len(unknownFeeds))
	})
}
