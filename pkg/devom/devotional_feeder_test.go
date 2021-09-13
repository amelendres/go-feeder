package devom

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/amelendres/go-feeder/pkg/cloud"
	"github.com/amelendres/go-feeder/pkg/fs"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

var path = map[string]string{
	"feeds-ok":    "./_test_feeds-ok.docx",
	"feeds-ko":    "./_test_feeds-ko.docx",
	"no-file":     "./_test_not-exist-file.docx",
	"drive-2019a": "1frfbhH2oUVOHLK7aNWr-0-2--hemIccj",
	"drive-2019b": "1jmcatkzNedm1aT9Y3JHjB51V4MdmoAS2",
	"drive-2019c": "1SQjizNJdE1QaIpbMB6oMcwHjc8_Lajue",
}

func TestDevotionalFeeder_FS(t *testing.T) {

	fp := fs.FileProvider{}
	dp := DevotionalParser{}
	df := NewFeeder(&fp, &dp)

	t.Run("it reads Feeds with UnknownFeeds from Docx", func(t *testing.T) {
		feeds, err := df.Feeds(path["feeds-ko"])

		assert.Nil(t, err)
		assert.Equal(t, 5, len(feeds.Feeds))
		assert.Equal(t, 5, len(feeds.UnknownFeeds))
	})

	t.Run("it fails read feeds without resource file", func(t *testing.T) {
		feeds, err := df.Feeds(path["no-file"])

		assert.NotNil(t, err)
		assert.Nil(t, feeds)
	})

	t.Run("it reads Feeds from Docx", func(t *testing.T) {
		feeds, err := df.Feeds(path["feeds-ok"])

		assert.Empty(t, err)
		assert.Empty(t, feeds.UnknownFeeds)
		assert.Equal(t, 14, len(feeds.Feeds))
	})
}

func TestDevotionalFeeder_GD(t *testing.T) {
	googleAPIKey := os.Getenv("GOOGLE_API_KEY")
	if googleAPIKey == "" {
		log.Fatal("ERROR: you must provide a Google Api Key")
	}

	ctx := context.Background()
	driveService, _ := drive.NewService(ctx, option.WithAPIKey(googleAPIKey))

	fp := cloud.NewGDFileProvider(driveService)
	dp := DevotionalParser{}
	df := NewFeeder(fp, &dp)

	t.Run("it reads from Google Drive", func(t *testing.T) {
		feeds, err := df.Feeds(path["drive-2019c"])

		assert.Nil(t, err)
		assert.Equal(t, 100, len(feeds.Feeds))
		assert.Equal(t, 0, len(feeds.UnknownFeeds))
	})
}
