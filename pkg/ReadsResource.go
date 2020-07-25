package feeder

import "errors"

var ErrOpeningFile = errors.New("Error opening file")
var ErrReadingFile = errors.New("Error reading file")

type ReadsResource interface {
	Read(url string) (string, error)
}
