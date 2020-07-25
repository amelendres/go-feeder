package feeder

import "errors"

var ErrUnknownFeed = errors.New("unknown feed")

type ReadsFeed interface {
	Feeds(path string) ([]Feed, []UnknownFeed, error)
}
