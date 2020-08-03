package feeder

type Feed []string

type UnknownFeed struct {
	Feed      []string
	FeedError string
}

type ParseFeeds struct {
	UnknownFeeds []UnknownFeed
	Feeds        []Feed
}
