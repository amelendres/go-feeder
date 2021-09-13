package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/amelendres/go-feeder/pkg/feeding"
	"github.com/amelendres/go-feeder/pkg/sending"

	"github.com/amelendres/go-feeder/pkg/cloud"
	"github.com/amelendres/go-feeder/pkg/devom"
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
	parser := devom.NewDevotionalParser()
	feeder := devom.NewDevotionalFeeder(fp, parser)
	sender := devom.NewPlanSender()

	ps := sending.NewService(sender, feeder)
	df := feeding.NewService(feeder)

	ds := server.NewFeederServer(ps, df)

	port := os.Getenv("PORT")
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), ds); err != nil {
		log.Fatalf("could not listen on port %s %v", port, err)
	}
}
