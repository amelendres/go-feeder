package server

import (
	"encoding/json"
	"github.com/amelendres/go-feeder/pkg/feeding"
	"github.com/amelendres/go-feeder/pkg/sending"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// DevServer is a HTTP interface for Cart
type DevServer struct {
	planSender sending.PlanSender
	devFeeder feeding.DevFeeder
	http.Handler
}

const jsonContentType = "application/json"

// NewDevServer creates a DevServer with routing configured
func NewDevServer(
	ps sending.PlanSender,
	df feeding.DevFeeder,
	) *DevServer {
	ds := &DevServer{planSender: ps, devFeeder: df}

	router := mux.NewRouter()
	router.Handle("/devotionals/import", http.HandlerFunc(ds.importDevotionalsHandler))
	router.Handle("/devotionals/parse", http.HandlerFunc(ds.parseDevotionalsHandler))

	ds.Handler = router

	return ds
}

func (ds *DevServer) importDevotionalsHandler(w http.ResponseWriter, r *http.Request) {

	var req sending.SendPlanReq
	json.NewDecoder(r.Body).Decode(&req)

	err := ds.planSender.Send(req)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (ds *DevServer) parseDevotionalsHandler(w http.ResponseWriter, r *http.Request) {

	var req feeding.FeedDevReq
	json.NewDecoder(r.Body).Decode(&req)

	feeds, err := ds.devFeeder.Feeds(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(feeds)

	w.WriteHeader(http.StatusOK)
}

