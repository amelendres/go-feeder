package devom_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/amelendres/go-feeder/internal/devom"
	feed "github.com/amelendres/go-feeder/pkg"
	"github.com/amelendres/go-feeder/pkg/cloud"
	"github.com/amelendres/go-feeder/pkg/fs"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

var (
	devomAPIUrl = "localhost:8030"
	api         = *devom.NewAPI(devomAPIUrl)
	path        = map[string]string{
		"dev-ok":              "./_test_devotionals-ok.docx",
		"dev-ko":              "./_test_devotionals-ko.docx",
		"no-file":             "./_test_not-exist-file.docx",
		"drive-dev-2019a":     "https://docs.google.com/document/d/1XI0cxe6T1VSipeeCmEbk14VDkZM5PS_c/preview",
		"drive-dev-bad-title": "https://docs.google.com/document/d/1frfbhH2oUVOHLK7aNWr-0-2--hemIccj/preview",
		"drive-dev-2019b":     "https://docs.google.com/document/d/1jmcatkzNedm1aT9Y3JHjB51V4MdmoAS2/preview",
		"drive-dev-2019c":     "https://docs.google.com/document/d/1SQjizNJdE1QaIpbMB6oMcwHjc8_Lajue/preview",
	}
)

func TestDevotionalFeeder_FS(t *testing.T) {

	dp := devom.NewDevotionalParser(api)
	fp := &fs.FileProvider{}
	df := feed.NewFeeder(dp, []feed.FileProvider{fp})

	t.Run("it reads Feeds with UnknownFeeds", func(t *testing.T) {
		feeds, err := df.Feeds(path["dev-ko"])

		assert.Nil(t, err)
		assert.Equal(t, 4, len(feeds.Items))
		assert.Equal(t, 6, len(feeds.UnknownItems))
	})

	t.Run("it fails read feeds without resource file", func(t *testing.T) {
		feeds, err := df.Feeds(path["no-file"])

		assert.NotNil(t, err)
		assert.Nil(t, feeds)
	})

	t.Run("it reads valid Feeds", func(t *testing.T) {
		feeds, err := df.Feeds(path["dev-ok"])

		assert.Empty(t, err)
		assert.Empty(t, feeds.UnknownItems)
		assert.Equal(t, 15, len(feeds.Items))
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
	dp := devom.NewDevotionalParser(api)
	df := feed.NewFeeder(dp, []feed.FileProvider{fp})

	t.Run("it reads from Google Drive", func(t *testing.T) {
		feeds, err := df.Feeds(path["drive-dev-2019a"])

		assert.Nil(t, err)
		assert.Equal(t, 100, len(feeds.Items))
		assert.Equal(t, 0, len(feeds.UnknownItems))
	})
}
