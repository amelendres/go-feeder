package main

import (
	"log"
	"net/http"

	"github.com/amelendres/go-feeder/pkg/devom"
	"github.com/amelendres/go-feeder/pkg/fs"
	"github.com/amelendres/go-feeder/pkg/server"
)

const dbFileName = "cart.db.json"

func main() {
	fp := fs.LocalFileProvider{}
	parser := devom.DevotionalParser{}
	res := fs.NewDocResource(&fp)
	feeder := fs.NewDocFeeder(res, &parser)

	ds := server.NewDevServer(feeder)

	if err := http.ListenAndServe(":5500", ds); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
