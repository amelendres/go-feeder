package sending

import (
	feed "github.com/amelendres/go-feeder/pkg"
	"github.com/pkg/errors"
)

var ErrUnknownFeed = errors.New("Unknown feeds")

type PlanSender interface {
	Send(req SendPlanReq) error
}

type service struct {
	sender feed.Sender
	feeder feed.Feeder
}

func NewPlanSender(s feed.Sender, f feed.Feeder) PlanSender {
	return &service{sender: s, feeder: f}
}

func (ps *service) Send(req SendPlanReq) error {

	feeds, unknownFeeds, err := ps.feeder.Feeds(req.FileUrl)
	if err != nil {
		return err
	}

	if len(unknownFeeds) > 0 {
		return ErrUnknownFeed
	}

	return ps.sender.Send(feeds)
}
