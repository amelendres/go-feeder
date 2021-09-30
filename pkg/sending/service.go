package sending

import (
	feed "github.com/amelendres/go-feeder/pkg"
	"github.com/pkg/errors"
)

var ErrUnknownFeed = errors.New("Unknown feeds")

type Service interface {
	Send(req SendReq) error
}

type SendReq struct {
	PlanId, AuthorId, PublisherId, FileUrl string
}
type service struct {
	sender feed.Sender
	feeder feed.Feeder
}

func NewService(s feed.Sender, f feed.Feeder) Service {
	return &service{sender: s, feeder: f}
}

func (ps *service) Send(req SendReq) error {
	dest := feed.NewDestination(req.PlanId, req.PublisherId, req.AuthorId)
	ps.feeder.Destination(dest)
	feeds, err := ps.feeder.Feeds(req.FileUrl)
	if err != nil {
		return err
	}

	if len(feeds.UnknownItems) > 0 {
		return ErrUnknownFeed
	}

	ps.sender.Destination(dest)
	return ps.sender.Send(feeds.Items)
}
