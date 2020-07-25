package feeder

type FeedParser interface {
	Parse(txt string) ([]Feed, []UnknownFeed)
}
