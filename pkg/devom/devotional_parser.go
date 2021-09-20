package devom

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"code.sajari.com/docconv"
	feed "github.com/amelendres/go-feeder/pkg"
)

var (
	ErrFeedDoesNotHavePassage      = errors.New("Feed does not have passage")
	ErrFeedDoesNotHaveContent      = errors.New("Feed does not have content")
	ErrFeedDoesNotHaveValidPassage = errors.New("Feed does not have a valid passage")
	ErrReadingResource             = func(err error) error {
		return fmt.Errorf("Error reading document: %w", err)
	}
	ErrDoesNotHaveValidDay = func(want, got int) error {
		return fmt.Errorf("Invalid day: want %d, but got %d", want, got)
	}
	ErrTitleAlreadyExists = func(title string) error {
		return fmt.Errorf("Title \"%s\" already exists", title)
	}
)

type DevotionalParser struct {
	Items       map[string]*feed.Feed
	Devotionals map[string]*Devotional
}

func NewDevotionalParser() feed.Parser {
	return &DevotionalParser{}
}

func (dp *DevotionalParser) Parse(r io.Reader) (*feed.ParseFeeds, error) {
	feeds := []feed.Feed{}
	unknownFeeds := []feed.UnknownFeed{}

	txt, err := dp.read(r)
	if err != nil {
		// log.Println(fmt.Errorf("Error reading resource: %s, %v ", r, err))
		return &feed.ParseFeeds{unknownFeeds, feeds}, err
	}

	devs := splitDevotionals(txt)
	lastDay := 0
	for _, dev := range devs {
		f, err := parseDevotional(dev)
		if err != nil {
			//log.Println(err, dev)
			unknownFeeds = append(unknownFeeds, feed.UnknownFeed{lines(dev), err.Error()})
			continue
		}
		day, err := strconv.Atoi(f["day"])
		if err != nil {
			unknownFeeds = append(unknownFeeds, feed.UnknownFeed{lines(dev), err.Error()})
			continue
		}

		//validate days
		if len(unknownFeeds) == 0 {
			lastDay = day - 1
			if len(feeds) > 0 {
				lastDay, _ = strconv.Atoi(feeds[len(feeds)-1]["day"])
			}
			if day != lastDay+1 {
				//log.Println(err, dev)
				err = ErrDoesNotHaveValidDay(lastDay+1, day)
				unknownFeeds = append(unknownFeeds, feed.UnknownFeed{lines(dev), err.Error()})
				continue
			}
		}

		//validate title
		// dp.Items = feeds
		if err = dp.uniqueTitle(f["title"]); err != nil {
			unknownFeeds = append(unknownFeeds, feed.UnknownFeed{lines(dev), err.Error()})
			continue
		}

		feeds = append(feeds, f)
	}

	return &feed.ParseFeeds{unknownFeeds, feeds}, nil
}

func (dp *DevotionalParser) uniqueTitle(title string) error {

	_, ok := dp.Items[title]
	if ok {
		return ErrTitleAlreadyExists(title)
	}
	_, ok = dp.Devotionals[title]
	if ok {
		return ErrTitleAlreadyExists(title)
	}
	return nil
}

func (dp *DevotionalParser) read(r io.Reader) (string, error) {

	content, _, err := docconv.ConvertDoc(r)

	if err != nil {
		return "", ErrReadingResource(err)
	}

	return content, nil
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
	lines := lines(text)

	dev := make(map[string]string)
	dev["day"] = lines[0]
	dev["title"] = lines[titleIdx]

	if len(lines) < 4 {
		return nil, feed.ErrUnknownFeed
	}

	var bibleReadingIdx int
	dev["bible_reading"], bibleReadingIdx = bibleReading(lines)

	if bibleReadingIdx == titleIdx+1 {
		return nil, ErrFeedDoesNotHavePassage
	}

	if bibleReadingIdx > titleIdx+1 {
		passage, err := passage(lines, titleIdx+1, bibleReadingIdx-1)
		if err != nil {
			return nil, err
		}
		dev["passage_text"], dev["passage_reference"] = passage.Text, passage.Reference
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
		dev["passage_text"], dev["passage_reference"] = passage.Text, passage.Reference
		dev["content"] = content(lines, contentIdx, len(lines)-1)
	}

	return dev, nil
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
