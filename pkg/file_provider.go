package feed

import (
	"errors"
	"io"
)

var ErrUnknownFile = errors.New("unknown file")

type FileProvider interface {
	File(path string) (io.Reader, error)
	Name() string
}
