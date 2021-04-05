package feed

type Destination interface {

}

type Sender interface {
	Send(feeds []Feed) error
	Destination(info Destination)
}
