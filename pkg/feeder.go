package feed

import "errors"

var ErrUnknownFeed = errors.New("unknown feed")

type Feeder interface {
	Feeds(path string) ([]Feed, []UnknownFeed, error)
}
