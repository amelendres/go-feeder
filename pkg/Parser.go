package feed

type Parser interface {
	Parse(txt string) ([]Feed, []UnknownFeed)
}
