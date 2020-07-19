package fs

import (
	"log"
	"regexp"
	"strings"

	feeder "github.com/amelendres/go-feeder/pkg"
)

type DocFeeder struct {
	resource feeder.ReadsResource
	feeds    []feeder.Feed
}

func NewDocFeeder(r feeder.ReadsResource) *DocFeeder {

	return &DocFeeder{
		resource: r,
	}
}

func (dr *DocFeeder) Feeds(path string) ([]feeder.Feed, error) {
	text, err := dr.resource.Read(path)

	//fmt.Println(text)
	feeds := parse(text)

	if err != nil {
		log.Fatal(err)
	}
	return feeds, nil
}

func parse(text string) []feeder.Feed {
	feeds := []feeder.Feed{}
	devs := splitDevotionals(text)

	for _, dev := range devs {
		devFeed := parseDevotional(dev)
		feeds = append(feeds, devFeed)
	}
	//fmt.Println(feeds)
	return feeds
}

func splitDevotionals(text string) []string {
	day := regexp.MustCompile(`\n\d{3}`)
	//todo: add day
	devs := day.Split(text, -1)
	devs = trimSlice(devs)
	return devs
}

func parseDevotional(text string) feeder.Feed {

	lines := strings.Split(text, "\n")
	lines = trimSlice(lines)
	var feed []string

	feed = append(feed, "###", lines[0], lines[1])
	var contentIdx int
	if isBibleReading(lines[2]) {
		feed = append(feed, lines[2])
		contentIdx = 4
	} else {
		feed = append(feed, "")
		contentIdx = 3
	}

	var content string
	for i := contentIdx; i < len(lines); i++ {
		content += lines[i]

	}
	feed = append(feed, content)

	// fmt.Println(feed)

	return feed
}

func isBibleReading(txt string) bool {
	return strings.Contains(txt, "Lectura:")
}

func trimSlice(slice []string) []string {
	var newSlice []string
	for _, item := range slice {
		if strings.TrimSpace(item) != "" {
			newSlice = append(newSlice, item)
		}
	}
	return newSlice
}
