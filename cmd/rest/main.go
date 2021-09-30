package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	feed "github.com/amelendres/go-feeder/pkg"
	"github.com/amelendres/go-feeder/pkg/feeding"
	"github.com/amelendres/go-feeder/pkg/fs"
	"github.com/amelendres/go-feeder/pkg/sending"

	"github.com/amelendres/go-feeder/internal/devom"
	"github.com/amelendres/go-feeder/pkg/cloud"
	"github.com/amelendres/go-feeder/pkg/server"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const (
	googleAPIKey = ""
	devomAPIUrl  = "http://localhost:8030/api/v1"
	serverPort   = "5500"
)

func main() {
	var (
		googleAPIKey = getEnv("GOOGLE_API_KEY", googleAPIKey)
		devomAPIUrl  = getEnv("DEVOM_API_URL", devomAPIUrl)
		serverPort   = getEnv("PORT", serverPort)
	)

	if googleAPIKey == "" {
		log.Fatal("ERROR: you must provide a Google Api Key")
	}

	ctx := context.Background()
	driveService, err := drive.NewService(ctx, option.WithAPIKey(googleAPIKey))
	if err != nil {
		log.Fatal("Unable start Drive Service")
	}

	//TODO switch port DEVOTIONAL==5500|TOPIC==5501
	gdp := cloud.NewGDFileProvider(driveService)
	fsp := fs.NewFileProvider()
	fileProviders := []feed.FileProvider{fsp, gdp}

	api := *devom.NewAPI(devomAPIUrl)
	parser := devom.NewTopicParser(api)
	feeder := feed.NewFeeder(parser, fileProviders)
	sender := devom.NewTopicSender(api)

	ps := sending.NewService(sender, feeder)
	df := feeding.NewService(feeder)

	ds := server.NewFeederServer(ps, df)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", serverPort), ds); err != nil {
		log.Fatalf("could not listen on port %s %v", serverPort, err)
	}
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}
