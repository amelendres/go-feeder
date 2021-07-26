package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/amelendres/go-feeder/pkg/cloud"
	"github.com/amelendres/go-feeder/pkg/devom"
	"github.com/amelendres/go-feeder/pkg/feeding"
	"github.com/amelendres/go-feeder/pkg/fs"
	"github.com/amelendres/go-feeder/pkg/sending"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	feeder "github.com/amelendres/go-feeder/pkg"
	"github.com/google/uuid"
)

var path = map[string]string{
	"feeds-ok":    "../fs/_test_feeds-ok.docx",
	"feeds-ko":    "../fs/_test_feeds-ko.docx",
	"feeds-2019a": "../fs/_test_feeds-2019a.docx",
	"feeds-2019b": "../fs/_test_feeds-2019b.docx",
	"no-file":     "../fs/_test_not-exist-file.docx",
	"drive-2019a": "1frfbhH2oUVOHLK7aNWr-0-2--hemIccj",
}

var planIds = map[int]string{
	2019: "1bec054b-ec6c-4ec0-becc-1a46bee429fb",
	2020: "a3f25740-d365-4c0e-8bd7-8dbb0a50cae3",
	2021: "23a63256-f264-4d94-b7ed-8ce60f744ae3",
	2022: "does-not-exists-this-yearly-plan-...",
}
var payload = sending.SendPlanReq{
	planIds[2021],
	uuid.New().String(),
	uuid.New().String(),
	"",
}

func TestImportDevotionals(t *testing.T) {

	fp := fs.NewFileProvider()
	parser := devom.NewParser()
	res := fs.NewDocResource(fp)
	feeder := fs.NewDocFeeder(res, parser)
	sender := devom.NewPlanSender()

	ps := sending.NewPlanSender(sender, feeder)
	df := feeding.NewDevFeeder(feeder)

	ds := NewDevServer(ps, df)

	t.Run("Unknown feeds", func(t *testing.T) {
		payload.FileUrl = path["feeds-ko"]

		response := httptest.NewRecorder()
		ds.ServeHTTP(response, newPostImportDevotionalRequest(payload))
		assert.Equal(t, http.StatusConflict, response.Code)
	})

	t.Run("Non existing resource file", func(t *testing.T) {
		payload.PlanId = planIds[2019]
		payload.FileUrl = path["no-file"]

		response := httptest.NewRecorder()
		ds.ServeHTTP(response, newPostImportDevotionalRequest(payload))
		assert.Equal(t, http.StatusConflict, response.Code)
	})

	t.Run("A new devotional feeds", func(t *testing.T) {
		payload.FileUrl = path["feeds-ok"]

		response := httptest.NewRecorder()
		ds.ServeHTTP(response, newPostImportDevotionalRequest(payload))
		assert.Equal(t, http.StatusAccepted, response.Code)
	})
}

func TestParseDevotionals(t *testing.T) {

	fp := fs.NewFileProvider()
	parser := devom.NewParser()
	res := fs.NewDocResource(fp)
	feeder := fs.NewDocFeeder(res, parser)
	sender := devom.NewPlanSender()

	ps := sending.NewPlanSender(sender, feeder)
	df := feeding.NewDevFeeder(feeder)

	ds := NewDevServer(ps, df)

	t.Run("A valid devotional feeds", func(t *testing.T) {
		payload.FileUrl = path["feeds-ok"]

		response := httptest.NewRecorder()

		ds.ServeHTTP(response, newPostParseDevotionalRequest(payload))

		parseFeeds := getParseFeedsFromResponse(t, response.Body)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, 14, len(parseFeeds.Feeds))
		assert.Equal(t, 0, len(parseFeeds.UnknownFeeds))
	})

	t.Run("With unknown devotional feeds", func(t *testing.T) {
		payload.FileUrl = path["feeds-ko"]
		response := httptest.NewRecorder()

		ds.ServeHTTP(response, newPostParseDevotionalRequest(payload))

		parseFeeds := getParseFeedsFromResponse(t, response.Body)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, 5, len(parseFeeds.Feeds))
		assert.Equal(t, 5, len(parseFeeds.UnknownFeeds))
	})

	t.Run("A non existing resource file", func(t *testing.T) {
		payload.FileUrl = path["no-file"]
		response := httptest.NewRecorder()

		ds.ServeHTTP(response, newPostParseDevotionalRequest(payload))

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
}

func TestGParseDevotionals(t *testing.T) {

	googleAPIKey := os.Getenv("GOOGLE_API_KEY")
	if googleAPIKey == "" {
		log.Fatal("ERROR: you must provide a Google Api Key")
	}
	ctx := context.Background()
	driveService, _ := drive.NewService(ctx, option.WithAPIKey(googleAPIKey))

	fp := cloud.NewGDFileProvider(driveService)
	parser := devom.NewParser()
	res := fs.NewDocResource(fp)
	feeder := fs.NewDocFeeder(res, parser)
	sender := devom.NewPlanSender()

	ps := sending.NewPlanSender(sender, feeder)
	df := feeding.NewDevFeeder(feeder)

	ds := NewDevServer(ps, df)

	t.Run("it parses from Google Drive File on POST", func(t *testing.T) {
		payload.FileUrl = path["drive-2019a"]

		response := httptest.NewRecorder()

		ds.ServeHTTP(response, newPostParseDevotionalRequest(payload))

		parseFeeds := getParseFeedsFromResponse(t, response.Body)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, 100, len(parseFeeds.Feeds))
		assert.Equal(t, 0, len(parseFeeds.UnknownFeeds))
	})
}

func newPostImportDevotionalRequest(sp sending.SendPlanReq) *http.Request {
	body, err := json.Marshal(sp)
	if err != nil {
		log.Fatalln(err)
	}
	req, _ := http.NewRequest(http.MethodPost, "/devotionals/import", bytes.NewBuffer(body))
	return req
}

func newPostParseDevotionalRequest(sp sending.SendPlanReq) *http.Request {
	body, err := json.Marshal(sp)
	if err != nil {
		log.Fatalln(err)
	}
	req, _ := http.NewRequest(http.MethodPost, "/devotionals/parse", bytes.NewBuffer(body))
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
