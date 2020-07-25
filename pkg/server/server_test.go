package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	feeder "github.com/amelendres/go-feeder/pkg"
	"github.com/amelendres/go-feeder/pkg/devom"
	"github.com/amelendres/go-feeder/pkg/fs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var path = map[string]string{
	"feeds-10-0": "../fs/_test_feeds-10-0.docx",
	"feeds-8-2":  "../fs/_test_feeds-8-2.docx",
	"no-file":    "../fs/_test_not-exist-file.docx",
	"2019a":      "../fs/_test_2019a.docx",
}

var planIds = map[int]string{
	2019: "f889d0c9-baf9-483f-8cfe-da1c5e0e3572",
	2020: "a3f25740-d365-4c0e-8bd7-8dbb0a50cae3",
	2021: "23a63256-f264-4d94-b7ed-8ce60f744ae3",
	2022: "does-not-exists-this-yearly-plan-...",
}
var payload = devom.ImportDailyDevotionals{
	planIds[2021],
	uuid.New().String(),
	uuid.New().String(),
	"",
}

func TestImportDevotionals(t *testing.T) {

	fp := fs.LocalFileProvider{}
	parser := devom.DevotionalParser{}
	res := fs.NewDocResource(&fp)
	feeder := fs.NewDocFeeder(res, &parser)

	ds := NewDevServer(feeder)

	t.Run("it import Devotionals on POST", func(t *testing.T) {
		payload.FileUrl = path["feeds-10-0"]

		response := httptest.NewRecorder()
		ds.ServeHTTP(response, newPostImportDevotionalRequest(payload))
		assert.Equal(t, http.StatusAccepted, response.Code)
		// assert.Equal(t, 10, len(getDevotionalsFromResponse(t, response.Body)))
		//assert new devotionals
	})

	t.Run("it fails import Devotionals on POST", func(t *testing.T) {
		payload.FileUrl = path["feeds-8-2"]

		response := httptest.NewRecorder()
		ds.ServeHTTP(response, newPostImportDevotionalRequest(payload))
		assert.Equal(t, http.StatusConflict, response.Code)
		// assert.Equal(t, 10, len(getDevotionalsFromResponse(t, response.Body)))
		//assert new devotionals
	})

	t.Run("it fails import not exists file on POST", func(t *testing.T) {
		payload.PlanId = planIds[2019]
		payload.FileUrl = path["no-file"]

		response := httptest.NewRecorder()
		ds.ServeHTTP(response, newPostImportDevotionalRequest(payload))
		assert.Equal(t, http.StatusBadRequest, response.Code)
		// assert.Equal(t, 10, len(getDevotionalsFromResponse(t, response.Body)))
		//assert new devotionals
	})

}

func TestParseDevotionals(t *testing.T) {

	fp := fs.LocalFileProvider{}
	parser := devom.DevotionalParser{}
	res := fs.NewDocResource(&fp)
	feeder := fs.NewDocFeeder(res, &parser)

	ds := NewDevServer(feeder)

	t.Run("it parses 10 Feeds on POST", func(t *testing.T) {
		payload.FileUrl = path["feeds-10-0"]

		response := httptest.NewRecorder()

		ds.ServeHTTP(response, newPostParseDevotionalRequest(payload))

		parseFeeds := getParseFeedsFromResponse(t, response.Body)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, 10, len(parseFeeds.Feeds))
		assert.Equal(t, 0, len(parseFeeds.UnknownFeeds))
	})

	t.Run("it parses 8 Feeds and 2 UnknowFeed on POST", func(t *testing.T) {
		payload.FileUrl = path["feeds-8-2"]
		response := httptest.NewRecorder()

		ds.ServeHTTP(response, newPostParseDevotionalRequest(payload))

		parseFeeds := getParseFeedsFromResponse(t, response.Body)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, 8, len(parseFeeds.Feeds))
		assert.Equal(t, 2, len(parseFeeds.UnknownFeeds))
	})

	t.Run("it fails parse does not exists file on POST", func(t *testing.T) {
		payload.FileUrl = path["no-file"]
		response := httptest.NewRecorder()

		ds.ServeHTTP(response, newPostParseDevotionalRequest(payload))

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
}

func newPostImportDevotionalRequest(importDevotional devom.ImportDailyDevotionals) *http.Request {
	body, err := json.Marshal(importDevotional)
	if err != nil {
		log.Fatalln(err)
	}
	req, _ := http.NewRequest(http.MethodPost, "/devotionals/import", bytes.NewBuffer(body))
	return req
}

func newPostParseDevotionalRequest(importDevotional devom.ImportDailyDevotionals) *http.Request {
	body, err := json.Marshal(importDevotional)
	if err != nil {
		log.Fatalln(err)
	}
	req, _ := http.NewRequest(http.MethodPost, "/devotionals/parse", bytes.NewBuffer(body))
	return req
}

// func getDevotionalsFromResponse(t *testing.T, body io.Reader) []devom.Devotional {
// 	t.Helper()
// 	devs, err := newDevotionalFromJSON(body)

// 	if err != nil {
// 		t.Fatalf("Unable to parse response from server %q into slice of Devotional, '%v'", body, err)
// 	}

// 	return devs
// }

func getParseFeedsFromResponse(t *testing.T, body io.Reader) feeder.ParseFeeds {
	t.Helper()
	parsedFeeds, err := newParseFeedsFromJSON(body)

	if err != nil {
		t.Fatalf("Unable to parse response from server %q into ParseFeeds, '%v'", body, err)
	}

	return parsedFeeds
}

// func newDevotionalFromJSON(rdr io.Reader) ([]devom.Devotional, error) {
// 	var devs []devom.Devotional
// 	err := json.NewDecoder(rdr).Decode(&devs)

// 	if err != nil {
// 		err = fmt.Errorf("problem parsing Devotionals, %v", err)
// 	}

// 	return devs, err
// }

func newParseFeedsFromJSON(rdr io.Reader) (feeder.ParseFeeds, error) {
	var parseFeeds feeder.ParseFeeds
	err := json.NewDecoder(rdr).Decode(&parseFeeds)

	if err != nil {
		err = fmt.Errorf("problem parsing ParseFeeds, %v", err)
	}

	return parseFeeds, err
}
