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
	ErrUndefinedDestination        = errors.New("Undefined destination")
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

type devotionalParser struct {
	api         API
	to          *feed.Destination
	items       map[string]*feed.Item
	devotionals map[string]*Devotional
}

func NewDevotionalParser(api API) feed.Parser {
	return &devotionalParser{api: api}
}

func (dp *devotionalParser) Destination(d *feed.Destination) {
	dp.to = d
}

func (dp *devotionalParser) Parse(r io.Reader) (*feed.ParsedItems, error) {
	feeds := []feed.Item{}
	unknownFeeds := []feed.UnknownItem{}

	txt, err := dp.read(r)
	if err != nil {
		return &feed.ParsedItems{unknownFeeds, feeds}, err
	}

	_ = dp.refreshCache()

	dp.items = make(map[string]*feed.Item)
	devs := splitDevotionals(txt)
	lastDay := 0
	for _, dev := range devs {
		f, err := parseDevotional(dev)
		if err != nil {
			unknownFeeds = append(unknownFeeds, feed.UnknownItem{lines(dev), err.Error()})
			continue
		}
		day, err := strconv.Atoi(f["day"])
		if err != nil {
			unknownFeeds = append(unknownFeeds, feed.UnknownItem{lines(dev), err.Error()})
			continue
		}

		//validate sequencial days
		if len(unknownFeeds) == 0 {
			lastDay = day - 1
			if len(feeds) > 0 {
				lastDay, _ = strconv.Atoi(feeds[len(feeds)-1]["day"])
			}
			if day != lastDay+1 {
				err = ErrDoesNotHaveValidDay(lastDay+1, day)
				unknownFeeds = append(unknownFeeds, feed.UnknownItem{lines(dev), err.Error()})
				continue
			}
		}

		//validate title
		if err = dp.uniqueTitle(f["title"]); err != nil {
			unknownFeeds = append(unknownFeeds, feed.UnknownItem{lines(dev), err.Error()})
			continue
		}

		feeds = append(feeds, f)
		dp.items[f["title"]] = &f
	}

	return &feed.ParsedItems{unknownFeeds, feeds}, nil
}

func (dp *devotionalParser) uniqueTitle(title string) error {

	_, ok := dp.items[title]
	if ok {
		return ErrTitleAlreadyExists(title)
	}
	_, ok = dp.devotionals[title]
	if ok {
		return ErrTitleAlreadyExists(title)
	}
	return nil
}

func (dp *devotionalParser) read(r io.Reader) (string, error) {

	content, _, err := docconv.ConvertDoc(r)

	if err != nil {
		return "", ErrReadingResource(err)
	}

	return content, nil
}

func (dp *devotionalParser) refreshCache() error {
	dp.devotionals = make(map[string]*Devotional)

	if dp.to == nil {
		return nil
	}

	devotionals, err := dp.api.getDevotionals(dp.to.AuthorId)
	if err != nil {
		return err
	}

	for _, dev := range devotionals {
		dp.devotionals[dev.Title] = dev
	}
	return nil
}

func splitDevotionals(text string) []string {

	day := regexp.MustCompile(`\n([0-9]+)(\n|\s*\n)`)

	devTexts := day.Split(text, -1)
	devTexts = trimSlice(devTexts)
	days := day.FindAllString(text, -1)
	// fmt.Printf("text: %+v\n total: %d\n", text, len(devTexts))
	// fmt.Printf("days: %+v\n total: %d\n", days, len(days))
	var devs []string
	for i, item := range devTexts {
		devs = append(devs, strings.TrimSpace(days[i])+"\n"+item)
	}

	return devs
}

func parseDevotional(text string) (feed.Item, error) {
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
