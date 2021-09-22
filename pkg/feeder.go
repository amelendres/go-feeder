package feed

import "errors"

var ErrUnknownFeed = errors.New("unknown feed")

type Feeder interface {
	Feeds(path string) (*ParseFeeds, error)
	Destination(d Destination)
}
