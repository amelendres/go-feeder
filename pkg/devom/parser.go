package devom

import (
	"errors"
	feed "github.com/amelendres/go-feeder/pkg"
	"regexp"
	"strings"
)

var ErrFeedDoesNotHavePassage = errors.New("Feed does not have passage")
var ErrFeedDoesNotHaveContent = errors.New("Feed does not have content")
var ErrFeedDoesNotHaveValidPassage = errors.New("Feed does not have a valid passage")

type Parser struct{}

func NewParser() feed.Parser{
	return &Parser{}
}

func (dp *Parser) Parse(txt string) ([]feed.Feed, []feed.UnknownFeed) {
	feeds := []feed.Feed{}
	unknownFeeds := []feed.UnknownFeed{}

	devs := splitDevotionals(txt)
	for _, dev := range devs {
		f, err := parseDevotional(dev)
		if err != nil {
			//log.Println(err, dev)
			unknownFeeds = append(unknownFeeds, feed.UnknownFeed{lines(dev), err.Error()})
		} else {
			feeds = append(feeds, f)
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

func parseDevotional(text string) (feed.Feed, error) {
	titleIdx := 1

	dev := map[string]string{}

	lines := lines(text)
	dev["day"] = lines[0]
	dev["title"] = lines[titleIdx]

	if len(lines) < 4 {
		return nil, feed.ErrUnknownFeed
	}

	var bibleReadingIdx int
	dev["bibleReading"], bibleReadingIdx = bibleReading(lines)

	if bibleReadingIdx == titleIdx+1 {
		return nil, ErrFeedDoesNotHavePassage
	}

	if bibleReadingIdx > titleIdx+1 {
		passage, err := passage(lines, titleIdx+1, bibleReadingIdx-1)
		if err != nil {
			return nil, err
		}
		dev["passage.text"], dev["passage.reference"] = passage.Text, passage.Reference
		dev["content"] = content(lines, bibleReadingIdx+1, len(lines)-1)

	} else {
		contentIdx := contentIndex(lines)
		if contentIdx < 0 {
			// log.Println(ErrFeedDoesNotHaveContent, text)
			return nil, ErrFeedDoesNotHaveContent
		}

		if contentIdx == titleIdx+1 {
			return nil, ErrFeedDoesNotHavePassage
		}
		passage, err := passage(lines, titleIdx+1, contentIdx-1)
		if err != nil {
			return nil, err
		}
		dev["passage.text"], dev["passage.reference"] = passage.Text, passage.Reference
		dev["content"] = content(lines, contentIdx, len(lines)-1)
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
		content += lines[i] + "\n\n"
	}

	return content
}

func passage(lines []string, start int, end int) (Passage, error) {
	txt := lines[start]

	if start == end {
		text, ref, err := splitPassage(txt)
		return NewPassage(text, ref), err
	}

	var passage string
	for i := start; i <= end; i++ {
		if passage != "" {
			passage += "\n\n"
		}

		passage += lines[i]
	}
	return NewPassage(passage, ""), nil
}

func contentIndex(lines []string) int {
	index := 3
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

func splitPassage(txt string) (text string, reference string, err error) {
	var passage []string
	lastPassageChar := regexp.MustCompile(`(”|")(\s*)\(`)
	occurrences := lastPassageChar.FindAllString(txt, -1)

	if len(occurrences) == 0 {
		return txt, "", ErrFeedDoesNotHaveValidPassage
	}

	if len(occurrences) > 1 {
		return txt, "", nil
	} else {
		passage = lastPassageChar.Split(txt, -1)
		if len(passage) < 2 {
			return passage[0], "", ErrFeedDoesNotHaveValidPassage
		}
		passage[0] += `”`
		passage[1] = `(` + passage[1]
	}

	return passage[0], passage[1], nil
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
