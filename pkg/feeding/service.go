package feeding

import (
	feed "github.com/amelendres/go-feeder/pkg"
	"github.com/amelendres/go-feeder/pkg/devom"
)

type FeedReq struct {
	PlanId, AuthorId, PublisherId, FileUrl string
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

func (s *service) Feeds(req FeedReq) (*feed.ParseFeeds, error) {
	s.feeder.Destination(devom.NewDestination(req.PlanId, req.PublisherId, req.AuthorId))
	return s.feeder.Feeds(req.FileUrl)
}
