package feed

type Destination interface {

}

type Sender interface {
	Send(feeds []Feed) error
	//Send(feeds []Feed, to Destination) error
	Destination(info Destination)
	//Destination(info interface{})
}
