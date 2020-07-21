package fs

import (
	"regexp"
	"strings"

	feeder "github.com/amelendres/go-feeder/pkg"
)

type DevotionalParser struct{}

func (dp *DevotionalParser) Parse(txt string) []feeder.Feed {
	feeds := []feeder.Feed{}
	devs := splitDevotionals(txt)

	for _, dev := range devs {
		devFeed := parseDevotional(dev)
		feeds = append(feeds, devFeed)
	}
	// fmt.Println(feeds)
	return feeds
}

func splitDevotionals(text string) []string {
	day := regexp.MustCompile(`\n\d{3}`)

	devTexts := day.Split(text, -1)
	devTexts = trimSlice(devTexts)

	days := day.FindAllString(text, -1)

	var devs []string
	for i, item := range devTexts {
		devs = append(devs, days[i]+"\n"+item)
	}
	// fmt.Println(days)
	// fmt.Println(devs)

	return devs
}

func parseDevotional(text string) feeder.Feed {

	lines := strings.Split(text, "\n")
	lines = trimSlice(lines)
	var feed []string

	feed = append(feed, lines[0], lines[1], lines[2])
	contentIdx := 4
	if isBibleReading(lines[3]) {
		feed = append(feed, lines[3])
	} else {
		feed = append(feed, "")
		contentIdx = 3
	}

	var content string
	for i := contentIdx; i < len(lines); i++ {
		content += lines[i]
	}
	feed = append(feed, content)

	//fmt.Println(feed)

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
