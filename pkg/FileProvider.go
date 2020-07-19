package feeder

import (
	"errors"
	"os"
)

var ErrUnknownFile = errors.New("unknown file")

type FileProvider interface {
	File(path string) (*os.File, error)
}
