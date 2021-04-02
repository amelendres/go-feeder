package devom

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	feed "github.com/amelendres/go-feeder/pkg"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"strconv"
)
var (
	ErrAddingDailyDevotional = errors.New("does not add daily devotional")
	ErrSendingDevotional = errors.New("does not create devotional")
)

type Destination struct {
	PlanId      string
	PublisherId string
	AuthorId    string
}

func NewDestination(planId, publisherId, authorId string) feed.Destination {
	return Destination{planId, publisherId, authorId}
}

type PlanSender struct {
	to Destination
	apiUrl string
	//to feed.Destination
	//ApiUrl string
}

//func NewPlanSender(p Destination, ApiUrl string) feed.Sender{
//	return &PlanSender{to: p, ApiUrl: ApiUrl}
//}

func NewPlanSender() feed.Sender{
	return &PlanSender{}
}

//func (ps *PlanSender) Destination(info feed.Destination) {
func (ps *PlanSender) Destination(info feed.Destination) {
	ps.to = info.(Destination)
}

func (ps *PlanSender) Send(feeds []feed.Feed) error {
	ps.apiUrl = fmt.Sprintf("%s/devotionals", os.Getenv("DEVOM_API_URL"))

	for _, feed := range feeds {
		dev := ps.mapFeed(feed)
		if err := ps.sendDevotional(dev); err == nil {
			day, _ := strconv.Atoi(feed[0])
			err = ps.addDailyDevotional(DailyDevotional{day, dev.Id}, ps.to.PlanId)
			if err != nil {
				//log.Print(err)
				return err
			}
		}
	}
	return nil
}

func (ps *PlanSender) mapFeed(feed []string) Devotional {

	return Devotional{
		uuid.New().String(),
		feed[1],
		Passage{feed[2], feed[3]},
		feed[5],
		feed[4],
		nil,
		ps.to.AuthorId,
		ps.to.PublisherId,
		nil,
	}
}

func (ps *PlanSender) sendDevotional(dev Devotional) error {
	body, err := json.Marshal(dev)
	if err != nil {
		return err
	}

	resp, err := http.Post(ps.apiUrl, "json", bytes.NewBuffer(body))
	log.Printf("[%s] %s \n", "POST", ps.apiUrl)

	if err != nil {
		log.Printf("ERROR: [%s] %s \npayload: %s\n\n", "POST", ps.apiUrl, string(body))
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("STATUS ERROR: [%s] %s \npayload: %s\n\n reponse: status %d", "POST", ps.apiUrl, string(body), resp.StatusCode)
		return ErrSendingDevotional
	}

	return nil
}

func (ps *PlanSender) addDailyDevotional(dev DailyDevotional, planId string) error {
	body, err := json.Marshal(dev)
	if err != nil {
		return err
	}

	resp, err := http.Post(ps.apiUrl, "json", bytes.NewBuffer(body))
	log.Printf("[%s] %s \n\n", "POST", ps.apiUrl)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("ERROR: [%s] %s \npayload: %s\n\n", "POST", ps.apiUrl, string(body))
		return ErrAddingDailyDevotional
	}
	return nil
}
