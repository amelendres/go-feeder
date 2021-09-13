package devom

import (
	feed "github.com/amelendres/go-feeder/pkg"
)

type DevotionalFeeder struct {
	fileProvider feed.FileProvider
	parser       feed.Parser
	feeds        []feed.Feed
}

func NewDevotionalFeeder(fp feed.FileProvider, p feed.Parser) feed.Feeder {

	return &DevotionalFeeder{
		fileProvider: fp,
		parser:       p,
	}
}

func (df *DevotionalFeeder) Feeds(path string) (*feed.ParseFeeds, error) {
	f, err := df.fileProvider.File(path)
	if err != nil {
		return nil, err
	}

	return df.parser.Parse(f)
}
