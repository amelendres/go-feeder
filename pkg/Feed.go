package feeder

type Feed []string
type UnknownFeed []string

type ParseFeeds struct {
	Feeds        []Feed
	UnknownFeeds []UnknownFeed
}
