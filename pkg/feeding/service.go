package feeding

import (
	feed "github.com/amelendres/go-feeder/pkg"
)

type FeedReq struct {
	PlanId, AuthorId, PublisherId, FileUrl string
}

type Service interface {
	Feeds(req FeedReq) (*feed.ParsedItems, error)
}

type service struct {
	feeder feed.Feeder
}

func NewService(f feed.Feeder) Service {
	return &service{feeder: f}
}

func (s *service) Feeds(req FeedReq) (*feed.ParsedItems, error) {
	s.feeder.Destination(feed.NewDestination(req.PlanId, req.PublisherId, req.AuthorId))
	return s.feeder.Feeds(req.FileUrl)
}
