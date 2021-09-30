package server

import (
	"encoding/json"
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
	router.Handle("/feeds/import", http.HandlerFunc(ds.importFeedHandler))
	router.Handle("/feeds/parse", http.HandlerFunc(ds.parseFeedHandler))

	ds.Handler = router

	return ds
}

func (ds *FeederServer) importFeedHandler(w http.ResponseWriter, r *http.Request) {

	var req sending.SendReq
	json.NewDecoder(r.Body).Decode(&req)

	err := ds.sender.Send(req)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (ds *FeederServer) parseFeedHandler(w http.ResponseWriter, r *http.Request) {

	var req feeding.FeedReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	feeds, err := ds.feeder.Feeds(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(feeds)
}
