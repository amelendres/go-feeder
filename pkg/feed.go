package feed

type Item map[string]string

type UnknownItem struct {
	Item      []string
	ItemError string
}

type ParsedItems struct {
	UnknownItems []UnknownItem
	Items        []Item
}
