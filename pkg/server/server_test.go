package server

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amelendres/go-feeder/pkg/devom"
	"github.com/amelendres/go-feeder/pkg/fs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestImportDevotional(t *testing.T) {

	importDevotional := devom.ImportDailyDevotionals{
		"23a63256-f264-4d94-b7ed-8ce60f744ae3",
		uuid.New().String(),
		uuid.New().String(),
		"../fs/Meditaciones 2019a.docx",
	}

	fp := fs.LocalFileProvider{}
	parser := devom.DevotionalParser{}
	res := fs.NewDocResource(&fp)
	feeder := fs.NewDocFeeder(res, &parser)

	ds := NewDevServer(feeder)

	t.Run("it import Devotionals on POST", func(t *testing.T) {
		response := httptest.NewRecorder()
		ds.ServeHTTP(response, newPostImportDevotionalRequest(importDevotional))
		assert.Equal(t, http.StatusAccepted, response.Code)

		//assert new devotionals
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
