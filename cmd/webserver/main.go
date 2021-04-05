package main

import (
	"context"
	"fmt"
	"github.com/amelendres/go-feeder/pkg/feeding"
	"github.com/amelendres/go-feeder/pkg/sending"
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
	parser := devom.NewParser()
	res := fs.NewDocResource(fp)
	feeder := fs.NewDocFeeder(res, parser)
	sender := devom.NewPlanSender()

	ps := sending.NewPlanSender(sender, feeder)
	df := feeding.NewDevFeeder(feeder)

	ds := server.NewDevServer(ps, df)

	port := os.Getenv("PORT")
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), ds); err != nil {
		log.Fatalf("could not listen on port %s %v", port, err)
	}
}
