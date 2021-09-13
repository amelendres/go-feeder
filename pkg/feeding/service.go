package feeding

import (
	feed "github.com/amelendres/go-feeder/pkg"
)

type FeedReq struct {
	FileUrl string
}

type Service interface {
	Feeds(req FeedReq) (*feed.ParseFeeds, error)
}

type service struct {
	feeder feed.Feeder
}

func NewService(f feed.Feeder) Service {
	return &service{feeder: f}
}

func (df *service) Feeds(req FeedReq) (*feed.ParseFeeds, error) {
	return df.feeder.Feeds(req.FileUrl)
}
