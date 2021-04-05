package feed

type Reader interface {
	Read(url string) (string, error)
}
