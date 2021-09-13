package sending

import (
	feed "github.com/amelendres/go-feeder/pkg"
	"github.com/amelendres/go-feeder/pkg/devom"
	"github.com/pkg/errors"
)

var ErrUnknownFeed = errors.New("Unknown feeds")

type Service interface {
	Send(req SendPlanReq) error
}

type service struct {
	sender feed.Sender
	feeder feed.Feeder
}

func NewService(s feed.Sender, f feed.Feeder) Service {
	return &service{sender: s, feeder: f}
}

func (ps *service) Send(req SendPlanReq) error {
	ps.sender.Destination(devom.NewDestination(req.PlanId, req.PublisherId, req.AuthorId))

	feeds, err := ps.feeder.Feeds(req.FileUrl)
	if err != nil {
		return err
	}

	if len(feeds.UnknownFeeds) > 0 {
		return ErrUnknownFeed
	}

	return ps.sender.Send(feeds.Feeds)
}
