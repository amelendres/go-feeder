package feed

//var ErrOpeningFile = errors.New("Error opening file: %v")
//var ErrReadingFile = errors.New("Error reading file")

type Reader interface {
	Read(url string) (string, error)
}
