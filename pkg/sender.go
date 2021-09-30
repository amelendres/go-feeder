package feed

type Destination struct {
	PlanId      string
	PublisherId string
	AuthorId    string
}

func NewDestination(planId, publisherId, authorId string) *Destination {
	return &Destination{planId, publisherId, authorId}
}

type Sender interface {
	Send(items []Item) error
	Destination(d *Destination)
}
