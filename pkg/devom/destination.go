package devom

import feed "github.com/amelendres/go-feeder/pkg"

type Destination struct {
	PlanId      string
	PublisherId string
	AuthorId    string
}

func NewDestination(planId, publisherId, authorId string) feed.Destination {
	return Destination{planId, publisherId, authorId}
}
