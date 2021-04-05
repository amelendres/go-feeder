package feeding

import (
	feed "github.com/amelendres/go-feeder/pkg"
)

type DevFeeder interface {
	Feeds(req FeedDevReq) (*feed.ParseFeeds, error)
}

type service struct {
	feeder feed.Feeder
}

func NewDevFeeder(f feed.Feeder) DevFeeder {
	return &service{feeder: f}
}

func (df *service) Feeds(req FeedDevReq) (*feed.ParseFeeds, error) {
	feeds, unknownFeeds, err := df.feeder.Feeds(req.FileUrl)
	if err != nil {
		return nil, err
	}
	return &feed.ParseFeeds{unknownFeeds, feeds}, nil
}

