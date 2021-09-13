package server

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

	"github.com/amelendres/go-feeder/pkg/cloud"
	"github.com/amelendres/go-feeder/pkg/devom"
	"github.com/amelendres/go-feeder/pkg/feeding"
	"github.com/amelendres/go-feeder/pkg/fs"
	"github.com/amelendres/go-feeder/pkg/sending"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"

	feeder "github.com/amelendres/go-feeder/pkg"
	"github.com/google/uuid"
)

var (
	path = map[string]string{
		"feeds-ok":    "../devom/_test_feeds-ok.docx",
		"feeds-ko":    "../devom/_test_feeds-ko.docx",
		"no-file":     "../devom/_test_not-exist-file.docx",
		"drive-2019a": "1frfbhH2oUVOHLK7aNWr-0-2--hemIccj",
	}

	planIds = map[int]string{
		2019: "1bec054b-ec6c-4ec0-becc-1a46bee429fb",
		2020: "a3f25740-d365-4c0e-8bd7-8dbb0a50cae3",
		2021: "23a63256-f264-4d94-b7ed-8ce60f744ae3",
		2022: "does-not-exists-this-yearly-plan-...",
	}
	payload = sending.SendReq{
		PlanId:      planIds[2021],
		AuthorId:    uuid.New().String(),
		PublisherId: uuid.New().String(),
		FileUrl:     "",
	}
)

// TODO: mock DEVOM server
func TestServer_ImportDevotionals(t *testing.T) {

	devomAPIUrl := os.Getenv("DEVOM_API_URL")
	fp := fs.NewFileProvider()
	parser := devom.NewDevotionalParser()
	feeder := devom.NewFeeder(fp, parser)
	sender := devom.NewPlanSender(devomAPIUrl)

	ps := sending.NewService(sender, feeder)
	df := feeding.NewService(feeder)

	ds := NewFeederServer(ps, df)

	t.Run("Unknown feeds", func(t *testing.T) {
		payload.FileUrl = path["feeds-ko"]

		response := httptest.NewRecorder()
		ds.ServeHTTP(response, newPostImportFeedRequest(payload))
		assert.Equal(t, http.StatusConflict, response.Code)
	})

	t.Run("Non existing resource file", func(t *testing.T) {
		payload.PlanId = planIds[2019]
		payload.FileUrl = path["no-file"]

		response := httptest.NewRecorder()
		ds.ServeHTTP(response, newPostImportFeedRequest(payload))
		assert.Equal(t, http.StatusConflict, response.Code)
	})

	t.Run("A new devotional feeds", func(t *testing.T) {
		payload.FileUrl = path["feeds-ok"]

		response := httptest.NewRecorder()
		ds.ServeHTTP(response, newPostImportFeedRequest(payload))
		assert.Equal(t, http.StatusAccepted, response.Code)
	})
}

func TestServer_ParseDevotionals(t *testing.T) {

	fp := fs.NewFileProvider()
	parser := devom.NewDevotionalParser()
	feeder := devom.NewFeeder(fp, parser)
	sender := devom.NewPlanSender("")

	ps := sending.NewService(sender, feeder)
	df := feeding.NewService(feeder)

	ds := NewFeederServer(ps, df)

	t.Run("A valid devotional feeds", func(t *testing.T) {
		payload.FileUrl = path["feeds-ok"]

		response := httptest.NewRecorder()

		ds.ServeHTTP(response, newPostParseFeedRequest(payload))

		parseFeeds := getParseFeedsFromResponse(t, response.Body)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, 14, len(parseFeeds.Feeds))
		assert.Equal(t, 0, len(parseFeeds.UnknownFeeds))
	})

	t.Run("With unknown devotional feeds", func(t *testing.T) {
		payload.FileUrl = path["feeds-ko"]
		response := httptest.NewRecorder()

		ds.ServeHTTP(response, newPostParseFeedRequest(payload))

		parseFeeds := getParseFeedsFromResponse(t, response.Body)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, 5, len(parseFeeds.Feeds))
		assert.Equal(t, 5, len(parseFeeds.UnknownFeeds))
	})

	t.Run("A non existing resource file", func(t *testing.T) {
		payload.FileUrl = path["no-file"]
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
	parser := devom.NewDevotionalParser()
	feeder := devom.NewFeeder(fp, parser)
	sender := devom.NewPlanSender("")

	ps := sending.NewService(sender, feeder)
	df := feeding.NewService(feeder)

	ds := NewFeederServer(ps, df)

	t.Run("it parses from Google Drive File on POST", func(t *testing.T) {
		payload.FileUrl = path["drive-2019a"]

		response := httptest.NewRecorder()

		ds.ServeHTTP(response, newPostParseFeedRequest(payload))

		parseFeeds := getParseFeedsFromResponse(t, response.Body)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, 100, len(parseFeeds.Feeds))
		assert.Equal(t, 0, len(parseFeeds.UnknownFeeds))
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

func getParseFeedsFromResponse(t *testing.T, body io.Reader) feeder.ParseFeeds {
	t.Helper()
	parsedFeeds, err := newParseFeedsFromJSON(body)

	if err != nil {
		t.Fatalf("Unable to parse response from server %q into ParseFeeds, '%v'", body, err)
	}

	return parsedFeeds
}

func newParseFeedsFromJSON(rdr io.Reader) (feeder.ParseFeeds, error) {
	var parseFeeds feeder.ParseFeeds
	err := json.NewDecoder(rdr).Decode(&parseFeeds)

	if err != nil {
		err = fmt.Errorf("problem parsing ParseFeeds, %v", err)
	}

	return parseFeeds, err
}
