package feeder

import "errors"

var ErrUnknownFeed = errors.New("unknown feed")

type FeedParser interface {
	Parse(txt string) ([]Feed, []UnknownFeed)
}
