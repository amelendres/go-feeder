package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	feeder "github.com/amelendres/go-feeder/pkg"
	"github.com/amelendres/go-feeder/pkg/devom"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// DevServer is a HTTP interface for Cart
type DevServer struct {
	feedReader feeder.ReadsFeed
	http.Handler
}

const jsonContentType = "application/json"

// NewDevServer creates a DevServer with routing configured
func NewDevServer(feedReader feeder.ReadsFeed) *DevServer {
	ds := new(DevServer)
	ds.feedReader = feedReader

	router := mux.NewRouter()
	router.Handle("/devotionals/import", http.HandlerFunc(ds.importDevotionalsHandler))

	ds.Handler = router

	return ds
}

func (ds *DevServer) importDevotionalsHandler(w http.ResponseWriter, r *http.Request) {

	var importData devom.ImportDailyDevotionals
	json.NewDecoder(r.Body).Decode(&importData)

	feeds, err := ds.feedReader.Feeds(importData.FileUrl)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, feed := range feeds {
		dev := buildDevotional(feed, importData)
		if devom.CreateDevotional(dev) {
			day, _ := strconv.Atoi(feed[0])
			devom.AddDailyDevotional(devom.DailyDevotional{day, dev.Id}, importData.PlanId)
		}

	}
	log.SetOutput(os.Stderr)

	w.WriteHeader(http.StatusAccepted)
}

func buildDevotional(feed []string, importData devom.ImportDailyDevotionals) devom.Devotional {
	passage := strings.Split(feed[2], "(")

	return devom.Devotional{
		uuid.New().String(),
		feed[1],
		devom.Passage{passage[0], "(" + passage[1]},
		feed[4],
		feed[3],
		nil,
		importData.AuthorId,
		importData.PublisherId,
		nil,
	}
}
