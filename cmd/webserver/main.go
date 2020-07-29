package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/amelendres/go-feeder/pkg/cloud"
	"github.com/amelendres/go-feeder/pkg/devom"
	"github.com/amelendres/go-feeder/pkg/fs"
	"github.com/amelendres/go-feeder/pkg/server"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const dbFileName = "cart.db.json"

func main() {
	googleAPIKey := os.Getenv("GOOGLE_API_KEY")
	if googleAPIKey == "" {
		log.Fatal("ERROR: you must provide a Google Api Key")
	}

	ctx := context.Background()
	driveService, err := drive.NewService(ctx, option.WithAPIKey(googleAPIKey))
	if err != nil {
		log.Fatal("Unable start Drive Service")
	}

	fp := cloud.NewGDFileProvider(driveService)
	parser := devom.DevotionalParser{}
	res := fs.NewDocResource(fp)
	feeder := fs.NewDocFeeder(res, &parser)

	ds := server.NewDevServer(feeder)

	if err := http.ListenAndServe(":5500", ds); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
