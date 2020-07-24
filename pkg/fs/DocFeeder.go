package fs

import (
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

func (dr *DocFeeder) Feeds(path string) ([]feeder.Feed, []feeder.UnknownFeed, error) {
	text, err := dr.resource.Read(path)
	if err != nil {
		// log.Fatal(err)
		return nil, nil, err
	}

	feeds, unknownFeeds := dr.parser.Parse(text)

	return feeds, unknownFeeds, nil
}
