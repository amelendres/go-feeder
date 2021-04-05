package fs

import (
	feed "github.com/amelendres/go-feeder/pkg"
)

type DocFeeder struct {
	resource feed.Reader
	parser   feed.Parser
	feeds    []feed.Feed
}

func NewDocFeeder(r feed.Reader, p feed.Parser) feed.Feeder {

	return &DocFeeder{
		resource: r,
		parser:   p,
	}
}

func (dr *DocFeeder) Feeds(path string) ([]feed.Feed, []feed.UnknownFeed, error) {
	text, err := dr.resource.Read(path)
	if err != nil {
		return nil, nil, err
	}

	feeds, unknownFeeds := dr.parser.Parse(text)

	return feeds, unknownFeeds, nil
}
