package devom

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"

	feed "github.com/amelendres/go-feeder/pkg"
)

const (
	topicDevotionalLength = 9
)

var (
	ErrInvalidDevotionalCell = func(text string) error {
		return fmt.Errorf("Invalid devotional cell <%s>", text)
	}
	ErrInvalidYear = func(year string) error {
		return fmt.Errorf("Invalid year <%s>", year)
	}
	ErrInvalidDay = func(day string) error {
		return fmt.Errorf("Invalid devotional <%s>", day)
	}
)

type TopicParser struct {
	api API
	to  *feed.Destination
}

func NewTopicParser(api API) feed.Parser {
	return &TopicParser{api: api}
}

func (dp *TopicParser) Destination(d *feed.Destination) {
	dp.to = d
}

func (dp *TopicParser) Parse(r io.Reader) (*feed.ParsedItems, error) {
	feeds := []feed.Item{}
	unknownFeeds := []feed.UnknownItem{}

	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, ErrReadingResource(err)
	}

	rows, err := f.GetRows("Traspuesto")
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		item, err := parseFeedItem(row)
		if err != nil {
			unknownFeeds = append(unknownFeeds, feed.UnknownItem{row, err.Error()})
			continue
		}

		feeds = append(feeds, item)
	}
	return &feed.ParsedItems{unknownFeeds, feeds}, nil
}

func parseFeedItem(row []string) (feed.Item, error) {
	topic := make(map[string]string)
	var devotionals []*YearlyDevotional
	for idx, colCell := range row {
		if idx == 0 {
			topic["title"] = colCell
			continue
		}

		dev, err := parseYearlyDevotional(colCell)
		if err != nil {
			// TODO: parse all row and save errors in order to generate ONE ERROR with all failures
			// indicating the file coordenates
			return nil, err
		}

		devotionals = append(devotionals, dev)
	}
	jsonDevs, err := json.Marshal(devotionals)
	if err != nil {
		return nil, err
	}
	topic["devotionals"] = string(jsonDevs)
	return topic, nil
}

func parseYearlyDevotional(text string) (*YearlyDevotional, error) {
	if len(text) < topicDevotionalLength {
		return nil, ErrInvalidDevotionalCell(text)
	}
	dev := strings.Split(text, " ")
	year, err := strconv.Atoi(dev[0][0:4])
	if err != nil {
		return nil, errors.Wrap(err, ErrInvalidYear(dev[0][0:4]).Error())
	}
	day, err := strconv.Atoi(dev[1])
	if err != nil {
		return nil, errors.Wrap(err, ErrInvalidDay(dev[1]).Error())
	}
	return &YearlyDevotional{Year: year, Day: day}, nil
}
