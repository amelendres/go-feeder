package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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
	router.Handle("/devotionals/parse", http.HandlerFunc(ds.parseDevotionalsHandler))

	ds.Handler = router

	return ds
}

func (ds *DevServer) importDevotionalsHandler(w http.ResponseWriter, r *http.Request) {

	var importData devom.ImportDailyDevotionals
	json.NewDecoder(r.Body).Decode(&importData)

	feeds, unknownFeeds, err := ds.feedReader.Feeds(importData.FileUrl)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(unknownFeeds) > 0 {
		log.Print(feeder.ErrUnknownFeed, unknownFeeds)
		w.WriteHeader(http.StatusConflict)
		return
	}

	for _, feed := range feeds {
		dev := buildDevotional(feed, importData)
		if err = devom.CreateDevotional(dev); err == nil {
			day, _ := strconv.Atoi(feed[0])
			err = devom.AddDailyDevotional(devom.DailyDevotional{day, dev.Id}, importData.PlanId)
			if err != nil {
				log.Print(err)
				w.WriteHeader(http.StatusConflict)
				return
			}
		}
	}

	w.WriteHeader(http.StatusAccepted)
}

func (ds *DevServer) parseDevotionalsHandler(w http.ResponseWriter, r *http.Request) {

	var importData devom.ImportDailyDevotionals
	json.NewDecoder(r.Body).Decode(&importData)

	feeds, unknownFeeds, err := ds.feedReader.Feeds(importData.FileUrl)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", jsonContentType)
	parseFeeds := feeder.ParseFeeds{feeds, unknownFeeds}

	json.NewEncoder(w).Encode(parseFeeds)

	w.WriteHeader(http.StatusOK)
}

func buildDevotional(feed []string, importData devom.ImportDailyDevotionals) devom.Devotional {

	return devom.Devotional{
		uuid.New().String(),
		feed[1],
		devom.Passage{feed[2], feed[3]},
		feed[5],
		feed[4],
		nil,
		importData.AuthorId,
		importData.PublisherId,
		nil,
	}
}
