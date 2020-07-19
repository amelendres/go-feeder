package feeder

type ReadsResource interface {
	Read(url string) (string, error)
}
