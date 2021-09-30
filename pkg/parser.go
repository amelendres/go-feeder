package feed

import "io"

type Parser interface {
	Parse(r io.Reader) (*ParsedItems, error)
	Destination(d *Destination)
}
