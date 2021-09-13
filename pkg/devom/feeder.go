package devom

import (
	feed "github.com/amelendres/go-feeder/pkg"
)

type Feeder struct {
	fileProvider feed.FileProvider
	parser       feed.Parser
	feeds        []feed.Feed
}

func NewFeeder(fp feed.FileProvider, p feed.Parser) feed.Feeder {

	return &Feeder{
		fileProvider: fp,
		parser:       p,
	}
}

func (df *Feeder) Feeds(path string) (*feed.ParseFeeds, error) {
	f, err := df.fileProvider.File(path)
	if err != nil {
		return nil, err
	}

	return df.parser.Parse(f)
}
