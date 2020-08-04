package devom

import (
	"errors"
	"log"
	"regexp"
	"strings"

	feeder "github.com/amelendres/go-feeder/pkg"
)

var ErrFeedDoesNotHasPassage = errors.New("Feed doesn not has passage")
var ErrFeedDoesNotHasContent = errors.New("Feed doesn not has content")

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
	titleIdx := 1

	dev := map[string]string{}

	lines := lines(text)
	dev["day"] = lines[0]
	dev["title"] = lines[titleIdx]

	if len(lines) < 4 {
		log.Println(feeder.ErrUnknownFeed, "[X] FEED TOO SHORT", text)
		return nil, feeder.ErrUnknownFeed
	}

	var bibleReadingIdx int
	dev["bibleReading"], bibleReadingIdx = bibleReading(lines)

	if bibleReadingIdx == titleIdx+1 {
		log.Println(ErrFeedDoesNotHasPassage, "bibleReadingIdx: "+string(bibleReadingIdx), text)
		return nil, ErrFeedDoesNotHasPassage
	}

	if bibleReadingIdx > titleIdx+1 {
		dev["passage.text"], dev["passage.reference"] = passage(lines, titleIdx+1, bibleReadingIdx-1)
		dev["content"] = content(lines, bibleReadingIdx+1, len(lines)-1)
	} else {
		contentIdx := contentIndex(lines)
		if contentIdx < 0 {
			log.Println(ErrFeedDoesNotHasContent, text)
			return nil, ErrFeedDoesNotHasContent
		}

		if contentIdx == titleIdx+1 {
			log.Println(ErrFeedDoesNotHasPassage, text)
			return nil, ErrFeedDoesNotHasPassage
		}
		dev["passage.text"], dev["passage.reference"] = passage(lines, titleIdx+1, contentIdx-1)
		dev["content"] = content(lines, contentIdx, len(lines)-1)

		// log.Println(dev["day"], titleIdx, contentIdx, dev["passage.text"], dev["passage.reference"])

	}

	var feed []string
	feed = append(feed,
		dev["day"],
		dev["title"],
		dev["passage.text"],
		dev["passage.reference"],
		dev["bibleReading"],
		dev["content"])

	return feed, nil
}

func lines(txt string) []string {
	lines := strings.Split(txt, "\n")
	lines = trimSlice(lines)
	return lines
}

func content(lines []string, start int, end int) string {
	content := ""
	for i := start; i <= end; i++ {
		// if content != "" {
		// 	content += "\n"
		// }
		content += lines[i]
	}

	return content
}

func passage(lines []string, start int, end int) (text string, reference string) {
	txt := lines[start]

	if start == end {
		passage := splitPassage(txt)
		return passage[0], passage[1]

	}

	var passage string
	for i := start; i <= end; i++ {
		if passage != "" {
			passage += "\n\n"
		}

		passage += lines[i]
	}
	return passage, ""
}

func contentIndex(lines []string) int {
	index := -1
	for key, line := range lines {
		if isPassage(line) {
			index = key + 1
		} else {
			if index > 0 {
				return index
			}
		}
	}
	return index
}

func bibleReading(lines []string) (txt string, key int) {
	for key, line := range lines {
		if isBibleReading(line) {
			return line, key
		}
	}
	return "", -1
}

func isBibleReading(txt string) bool {
	return strings.Contains(txt, "Lectura:")
}

func isPassage(txt string) bool {
	txt = strings.TrimSpace(txt)
	match, _ := regexp.MatchString(`^[“|"](.*)[”|"](.*)\((.*)\).?$`, txt)
	return match
}

func splitPassage(txt string) []string {
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
