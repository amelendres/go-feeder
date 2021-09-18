package feed

type Destination interface{}

type Sender interface {
	Send(items []Feed, to Destination) error
}
