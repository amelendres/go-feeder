package devom

import (
	"errors"
	"log"
	"regexp"
	"strings"

	feeder "github.com/amelendres/go-feeder/pkg"
)

var ErrFeedDoesNotHasPassage = errors.New("Feed doesn not has passage")

type DevotionalParser struct{}

func (dp *DevotionalParser) Parse(txt string) ([]feeder.Feed, []feeder.UnknownFeed) {
	feeds := []feeder.Feed{}
	unknownFeeds := []feeder.UnknownFeed{}

	devs := splitDevotionals(txt)
	for _, dev := range devs {
		feed, err := parseDevotional(dev)
		if err != nil {
			unknownFeeds = append(unknownFeeds, feeder.UnknownFeed{dev})
		} else {
			feeds = append(feeds, feed)
		}
	}
	return feeds, unknownFeeds
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

	return devs
}

func parseDevotional(text string) (feeder.Feed, error) {
	lines := lines(text)

	if len(lines) < 4 {
		return nil, feeder.ErrUnknownFeed
	}
	if !isPassage(lines[2]) {
		log.Println(ErrFeedDoesNotHasPassage, text)
		return nil, ErrFeedDoesNotHasPassage
	}
	passage := splitPassage(lines[2])

	var feed []string
	feed = append(feed, lines[0], lines[1], passage[0], passage[1])
	contentIdx := 4
	if isBibleReading(lines[3]) {
		feed = append(feed, strings.Split(lines[3], "Lectura:")[1])
	} else {
		feed = append(feed, "")
		contentIdx = 3
	}

	var content string
	for i := contentIdx; i < len(lines); i++ {
		content += lines[i]
	}
	feed = append(feed, content)

	return feed, nil
}

func lines(txt string) []string {
	lines := strings.Split(txt, "\n")
	lines = trimSlice(lines)
	return lines
}

func isBibleReading(txt string) bool {
	return strings.Contains(txt, "Lectura:")
}

func isPassage(txt string) bool {
	txt = strings.TrimSpace(txt)
	match, _ := regexp.MatchString(`^[“|"](.*)[”|"](.*)\((.*)\)`, txt)
	return match
}

func splitPassage(txt string) []string {
	//todo: detect passage list and join into passage.text
	var passage []string

	lastPassageChar := regexp.MustCompile("”|\"\\s")

	if len(lastPassageChar.FindAllString(txt, -1)) > 1 {
		passage = append(passage, txt, "")
	} else {
		passage = lastPassageChar.Split(txt, -1)
		if len(passage) < 2 {
			log.Println(len(passage))
			log.Fatalln(passage)
		}
		passage[0] = passage[0] + `”`
	}

	return passage
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
