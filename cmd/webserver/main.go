package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/amelendres/go-feeder/pkg/feeding"
	"github.com/amelendres/go-feeder/pkg/fs"
	"github.com/amelendres/go-feeder/pkg/sending"

	"github.com/amelendres/go-feeder/pkg/devom"
	"github.com/amelendres/go-feeder/pkg/server"
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

	// ctx := context.Background()
	// driveService, err := drive.NewService(ctx, option.WithAPIKey(googleAPIKey))
	// if err != nil {
	// 	log.Fatal("Unable start Drive Service")
	// }
	// fp := cloud.NewGDFileProvider(driveService)
	fp := fs.NewFileProvider()

	api := *devom.NewAPI(devomAPIUrl)
	parser := devom.NewDevotionalParser(api)
	feeder := devom.NewFeeder(fp, parser)
	sender := devom.NewPlanSender(api)

	ps := sending.NewService(sender, feeder)
	df := feeding.NewService(feeder)

	ds := server.NewFeederServer(ps, df)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", serverPort), ds); err != nil {
		log.Fatalf("could not listen on port %s %v", serverPort, err)
	}
}

// TODO: refactor as an env package
func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}
