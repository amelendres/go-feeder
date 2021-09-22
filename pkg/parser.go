package feed

import "io"

type Parser interface {
	Parse(r io.Reader) (*ParseFeeds, error)
	Destination(d Destination)
}
