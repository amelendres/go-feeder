package feeder

type ReadsFeed interface {
	Feeds(path string) ([]Feed, []UnknownFeed, error)
}
