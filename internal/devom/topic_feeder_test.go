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

var feedSource = map[string]string{
	"topics-ok":       "./_test_topics-ok.xlsx",
	"topics-ko":       "./_test_topics-ko.xlsx",
	"no-file":         "./_test_not-exist-file.xlsx",
	"drive-topics-ok": "https://docs.google.com/spreadsheets/d/1kWN7HHNrlytOyApwUlnA0SXc-WBsbuiA/preview",
}

func TestTopicFeeder_FS(t *testing.T) {

	fp := &fs.FileProvider{}
	dp := &devom.TopicParser{}
	df := feed.NewFeeder(dp, []feed.FileProvider{fp})

	t.Run("it parses Feed with UnknownFeeds", func(t *testing.T) {
		feeds, err := df.Feeds(feedSource["topics-ko"])

		assert.Nil(t, err)
		assert.Equal(t, 2, len(feeds.Items))
		assert.Equal(t, 5, len(feeds.UnknownItems))
	})

	t.Run("it fails read Feed without resource file", func(t *testing.T) {
		feeds, err := df.Feeds(feedSource["no-file"])

		assert.NotNil(t, err)
		assert.Nil(t, feeds)
	})

	t.Run("it parses Feed", func(t *testing.T) {
		feeds, err := df.Feeds(feedSource["topics-ok"])

		assert.Nil(t, err)
		assert.Equal(t, 7, len(feeds.Items))
		assert.Equal(t, 0, len(feeds.UnknownItems))
	})
}

func TestTopicFeeder_GD(t *testing.T) {
	googleAPIKey := os.Getenv("GOOGLE_API_KEY")
	if googleAPIKey == "" {
		log.Fatal("ERROR: you must provide a Google Api Key")
	}

	ctx := context.Background()
	driveService, _ := drive.NewService(ctx, option.WithAPIKey(googleAPIKey))

	fp := cloud.NewGDFileProvider(driveService)
	dp := &devom.TopicParser{}
	df := feed.NewFeeder(dp, []feed.FileProvider{fp})

	t.Run("it parses from Google Drive", func(t *testing.T) {
		feeds, err := df.Feeds(feedSource["drive-topics-ok"])

		assert.Nil(t, err)
		assert.Equal(t, 7, len(feeds.Items))
		assert.Equal(t, 0, len(feeds.UnknownItems))
	})
}
