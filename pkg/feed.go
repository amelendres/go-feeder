package feed

//TODO: Rename to Item
type Feed map[string]string

//TODO: Rename to UnknownItem
type UnknownFeed struct {
	Feed      []string
	FeedError string
}

//TODO: Rename to ParsedItems
type ParseFeeds struct {
	UnknownFeeds []UnknownFeed
	Feeds        []Feed
}
