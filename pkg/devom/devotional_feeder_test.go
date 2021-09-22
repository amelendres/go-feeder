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

var (
	devomAPIUrl = "localhost:8030"
	api         = *NewAPI(devomAPIUrl)
	path        = map[string]string{
		"dev-ok":              "./_test_devotionals-ok.docx",
		"dev-ko":              "./_test_devotionals-ko.docx",
		"no-file":             "./_test_not-exist-file.docx",
		"drive-dev-2019a":     "1XI0cxe6T1VSipeeCmEbk14VDkZM5PS_c",
		"drive-dev-bad-title": "1frfbhH2oUVOHLK7aNWr-0-2--hemIccj",
		"drive-dev-2019b":     "1jmcatkzNedm1aT9Y3JHjB51V4MdmoAS2",
		"drive-dev-2019c":     "1SQjizNJdE1QaIpbMB6oMcwHjc8_Lajue",
	}
)

func TestDevotionalFeeder_FS(t *testing.T) {

	dp := NewDevotionalParser(api)
	fp := &fs.FileProvider{}
	df := NewFeeder(fp, dp)

	t.Run("it reads Feeds with UnknownFeeds", func(t *testing.T) {
		feeds, err := df.Feeds(path["dev-ko"])

		assert.Nil(t, err)
		assert.Equal(t, 4, len(feeds.Feeds))
		assert.Equal(t, 6, len(feeds.UnknownFeeds))
	})

	t.Run("it fails read feeds without resource file", func(t *testing.T) {
		feeds, err := df.Feeds(path["no-file"])

		assert.NotNil(t, err)
		assert.Nil(t, feeds)
	})

	t.Run("it reads valid Feeds", func(t *testing.T) {
		feeds, err := df.Feeds(path["dev-ok"])

		assert.Empty(t, err)
		assert.Empty(t, feeds.UnknownFeeds)
		assert.Equal(t, 15, len(feeds.Feeds))
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
	dp := NewDevotionalParser(api)
	df := NewFeeder(fp, dp)

	t.Run("it reads from Google Drive", func(t *testing.T) {
		feeds, err := df.Feeds(path["drive-dev-2019a"])

		assert.Nil(t, err)
		assert.Equal(t, 100, len(feeds.Feeds))
		assert.Equal(t, 0, len(feeds.UnknownFeeds))
	})
}
