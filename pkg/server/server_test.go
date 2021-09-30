package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/amelendres/go-feeder/internal/devom"
	feed "github.com/amelendres/go-feeder/pkg"
	"github.com/amelendres/go-feeder/pkg/cloud"
	"github.com/amelendres/go-feeder/pkg/feeding"
	"github.com/amelendres/go-feeder/pkg/fs"
	"github.com/amelendres/go-feeder/pkg/sending"
	"github.com/amelendres/go-feeder/pkg/server"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"

	feeder "github.com/amelendres/go-feeder/pkg"
)

var (
	devomAPIUrl = os.Getenv("DEVOM_API_URL")
	api         = *devom.NewAPI(devomAPIUrl)
	feedSource  = map[string]string{
		"dev-ok":              "../../internal/devom/_test_devotionals-ok.docx",
		"dev-ko":              "../../internal/devom/_test_devotionals-ko.docx",
		"no-file":             "../../internal/devom/_test_not-exists-file",
		"drive-dev-2019a":     "https://docs.google.com/document/d/1XI0cxe6T1VSipeeCmEbk14VDkZM5PS_c/preview",
		"drive-dev-bad-title": "https://docs.google.com/document/d/1frfbhH2oUVOHLK7aNWr-0-2--hemIccj/preview",
		"topics-ok":           "../../internal/devom/_test_topics-ok.xlsx",
		"topics-ko":           "../../internal/devom/_test_topics-ko.xlsx",
		"drive-topics-ok":     "https://docs.google.com/spreadsheets/d/1kWN7HHNrlytOyApwUlnA0SXc-WBsbuiA/preview",
	}

	planIds = map[int]string{
		2019: "1bec054b-ec6c-4ec0-becc-1a46bee429fb",
		2020: "a3f25740-d365-4c0e-8bd7-8dbb0a50cae3",
		2021: "23a63256-f264-4d94-b7ed-8ce60f744ae3",
		2022: "does-not-exists-this-yearly-plan-...",
	}
	payload = sending.SendReq{
		PlanId:      planIds[2021],
		AuthorId:    "eeef78ed-043e-40db-9eae-8a0d77950ceb",
		PublisherId: "e5f12936-1339-4dbc-b339-fb0ba42b13a9",
		FileUrl:     "",
	}
)

// TODO: mock DEVOM server
func TestServer_ImportDevotionals_FromFS(t *testing.T) {

	fp := fs.NewFileProvider()
	parser := devom.NewDevotionalParser(api)
	feeder := feed.NewFeeder(parser, []feed.FileProvider{fp})
	sender := devom.NewDevotionalSender(api)

	ps := sending.NewService(sender, feeder)
	df := feeding.NewService(feeder)

	ds := server.NewFeederServer(ps, df)

	t.Run("Unknown feeds", func(t *testing.T) {
		payload.FileUrl = feedSource["dev-ko"]

		response := httptest.NewRecorder()
		ds.ServeHTTP(response, newPostImportFeedRequest(payload))
		assert.Equal(t, http.StatusConflict, response.Code)
	})

	t.Run("Non existing resource file", func(t *testing.T) {
		payload.PlanId = planIds[2019]
		payload.FileUrl = feedSource["no-file"]

		response := httptest.NewRecorder()
		ds.ServeHTTP(response, newPostImportFeedRequest(payload))
		assert.Equal(t, http.StatusConflict, response.Code)
	})

	t.Run("A new devotional feeds", func(t *testing.T) {
		payload.FileUrl = feedSource["dev-ok"]

		response := httptest.NewRecorder()
		ds.ServeHTTP(response, newPostImportFeedRequest(payload))
		assert.Equal(t, http.StatusAccepted, response.Code)
	})
}

func TestServer_ImportTopics(t *testing.T) {

	fp := fs.NewFileProvider()
	parser := devom.NewTopicParser(api)
	feeder := feed.NewFeeder(parser, []feed.FileProvider{fp})
	sender := devom.NewTopicSender(api)

	ps := sending.NewService(sender, feeder)
	df := feeding.NewService(feeder)

	ds := server.NewFeederServer(ps, df)

	t.Run("Unknown items", func(t *testing.T) {
		payload.FileUrl = feedSource["topics-ko"]

		response := httptest.NewRecorder()
		ds.ServeHTTP(response, newPostImportFeedRequest(payload))
		assert.Equal(t, http.StatusConflict, response.Code)
	})

	t.Run("Non existing resource file", func(t *testing.T) {
		payload.PlanId = planIds[2019]
		payload.FileUrl = feedSource["no-file"]

		response := httptest.NewRecorder()
		ds.ServeHTTP(response, newPostImportFeedRequest(payload))
		assert.Equal(t, http.StatusConflict, response.Code)
	})

	t.Run("Valid feed items", func(t *testing.T) {
		payload.FileUrl = feedSource["topics-ok"]

		response := httptest.NewRecorder()
		ds.ServeHTTP(response, newPostImportFeedRequest(payload))
		assert.Equal(t, http.StatusAccepted, response.Code)
	})
}

func TestServer_ParseDevotionals_FromFS(t *testing.T) {

	fp := fs.NewFileProvider()
	parser := devom.NewDevotionalParser(api)
	feeder := feed.NewFeeder(parser, []feed.FileProvider{fp})
	sender := devom.NewDevotionalSender(api)

	ps := sending.NewService(sender, feeder)
	df := feeding.NewService(feeder)

	ds := server.NewFeederServer(ps, df)

	t.Run("A valid devotional feeds", func(t *testing.T) {
		payload.FileUrl = feedSource["dev-ok"]

		response := httptest.NewRecorder()

		ds.ServeHTTP(response, newPostParseFeedRequest(payload))

		parseFeeds := getParseFeedsFromResponse(t, response.Body)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, 15, len(parseFeeds.Items))
		assert.Equal(t, 0, len(parseFeeds.UnknownItems))
	})

	t.Run("With unknown devotional feeds", func(t *testing.T) {
		payload.FileUrl = feedSource["dev-ko"]
		response := httptest.NewRecorder()

		ds.ServeHTTP(response, newPostParseFeedRequest(payload))

		parseFeeds := getParseFeedsFromResponse(t, response.Body)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, 4, len(parseFeeds.Items))
		assert.Equal(t, 6, len(parseFeeds.UnknownItems))
	})

	t.Run("A non existing resource file", func(t *testing.T) {
		payload.FileUrl = feedSource["no-file"]
		response := httptest.NewRecorder()

		ds.ServeHTTP(response, newPostParseFeedRequest(payload))

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
}

func TestServer_ParseDevotionals_FromGoogleDrive(t *testing.T) {

	googleAPIKey := os.Getenv("GOOGLE_API_KEY")
	if googleAPIKey == "" {
		log.Fatal("ERROR: you must provide a Google Api Key")
	}
	ctx := context.Background()
	driveService, _ := drive.NewService(ctx, option.WithAPIKey(googleAPIKey))

	fp := cloud.NewGDFileProvider(driveService)
	parser := devom.NewDevotionalParser(api)
	feeder := feed.NewFeeder(parser, []feed.FileProvider{fp})
	sender := devom.NewDevotionalSender(api)

	ps := sending.NewService(sender, feeder)
	df := feeding.NewService(feeder)

	ds := server.NewFeederServer(ps, df)

	t.Run("it parses from Google Drive File on POST", func(t *testing.T) {
		payload.FileUrl = feedSource["drive-dev-2019a"]

		response := httptest.NewRecorder()

		ds.ServeHTTP(response, newPostParseFeedRequest(payload))

		parseFeeds := getParseFeedsFromResponse(t, response.Body)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, 100, len(parseFeeds.Items))
		assert.Equal(t, 0, len(parseFeeds.UnknownItems))
	})
}

func newPostImportFeedRequest(sp sending.SendReq) *http.Request {
	body, err := json.Marshal(sp)
	if err != nil {
		log.Fatalln(err)
	}
	req, _ := http.NewRequest(http.MethodPost, "/feeds/import", bytes.NewBuffer(body))
	return req
}

func newPostParseFeedRequest(sp sending.SendReq) *http.Request {
	body, err := json.Marshal(sp)
	if err != nil {
		log.Fatalln(err)
	}
	req, _ := http.NewRequest(http.MethodPost, "/feeds/parse", bytes.NewBuffer(body))
	return req
}

func getParseFeedsFromResponse(t *testing.T, body io.Reader) feeder.ParsedItems {
	t.Helper()
	parsedFeeds, err := newParseFeedsFromJSON(body)

	if err != nil {
		t.Fatalf("Unable to parse response from server %q into ParseFeeds, '%v'", body, err)
	}

	return parsedFeeds
}

func newParseFeedsFromJSON(rdr io.Reader) (feeder.ParsedItems, error) {
	var parseFeeds feeder.ParsedItems
	err := json.NewDecoder(rdr).Decode(&parseFeeds)

	if err != nil {
		err = fmt.Errorf("problem parsing ParseFeeds, %v", err)
	}

	return parseFeeds, err
}
