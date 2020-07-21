package fs

import (
	"log"

	feeder "github.com/amelendres/go-feeder/pkg"
)

type DocFeeder struct {
	resource feeder.ReadsResource
	parser   feeder.FeedParser
	feeds    []feeder.Feed
}

func NewDocFeeder(r feeder.ReadsResource, p feeder.FeedParser) *DocFeeder {

	return &DocFeeder{
		resource: r,
		parser:   p,
	}
}

func (dr *DocFeeder) Feeds(path string) ([]feeder.Feed, error) {
	text, err := dr.resource.Read(path)

	//fmt.Println(text)
	feeds := dr.parser.Parse(text)

	if err != nil {
		log.Fatal(err)
	}
	return feeds, nil
}
