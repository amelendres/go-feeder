package feed

type Destination interface{}

type Sender interface {
	Send(feeds []Feed, to Destination) error
}
