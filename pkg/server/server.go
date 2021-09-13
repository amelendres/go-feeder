package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/amelendres/go-feeder/pkg/feeding"
	"github.com/amelendres/go-feeder/pkg/sending"
	"github.com/gorilla/mux"
)

type FeederServer struct {
	sender sending.Service
	feeder feeding.Service
	http.Handler
}

const jsonContentType = "application/json"

// NewFeederServer creates a DevServer with routing configured
func NewFeederServer(
	ss sending.Service,
	fs feeding.Service,
) *FeederServer {
	ds := &FeederServer{sender: ss, feeder: fs}

	router := mux.NewRouter()
	router.Handle("/devotionals/import", http.HandlerFunc(ds.importDevotionalsHandler))
	router.Handle("/devotionals/parse", http.HandlerFunc(ds.parseDevotionalsHandler))

	ds.Handler = router

	return ds
}

func (ds *FeederServer) importDevotionalsHandler(w http.ResponseWriter, r *http.Request) {

	var req sending.SendPlanReq
	json.NewDecoder(r.Body).Decode(&req)

	err := ds.sender.Send(req)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (ds *FeederServer) parseDevotionalsHandler(w http.ResponseWriter, r *http.Request) {

	var req feeding.FeedReq
	json.NewDecoder(r.Body).Decode(&req)

	feeds, err := ds.feeder.Feeds(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(feeds)

	w.WriteHeader(http.StatusOK)
}
