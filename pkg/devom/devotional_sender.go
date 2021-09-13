package devom

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	feed "github.com/amelendres/go-feeder/pkg"
	"github.com/google/uuid"
)

var (
	ErrAddingDailyDevotional = errors.New("does not add daily devotional")
	ErrSendingDevotional     = errors.New("does not create devotional")
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
	to     Destination
	apiUrl string
}

func NewPlanSender() feed.Sender {
	return &PlanSender{}
}

func (ps *PlanSender) Destination(info feed.Destination) {
	ps.to = info.(Destination)
}

func (ps *PlanSender) Send(feeds []feed.Feed) error {
	ps.apiUrl = os.Getenv("DEVOM_API_URL")

	for _, f := range feeds {
		dev := ps.mapFeed(f)
		if err := ps.sendDevotional(dev); err == nil {
			day, _ := strconv.Atoi(f[0])
			err = ps.addDailyDevotional(AddDailyDevotionalReq{day, dev.Id}, ps.to.PlanId)
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
	endpoint := fmt.Sprintf("%s/devotionals", ps.apiUrl)

	body, err := json.Marshal(dev)
	if err != nil {
		return err
	}

	resp, err := http.Post(endpoint, "json", bytes.NewBuffer(body))
	log.Printf("[%s] %s \n", "POST", endpoint)

	if err != nil {
		log.Printf("ERROR: [%s] %s \npayload: %s\n\n", "POST", endpoint, string(body))
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		log.Printf("STATUS ERROR: [%s] %s \npayload: %s\n\n reponse: status %d", "POST", endpoint, string(body), resp.StatusCode)
		return ErrSendingDevotional
	}

	return nil
}

func (ps *PlanSender) addDailyDevotional(req AddDailyDevotionalReq, planId string) error {
	endpoint := fmt.Sprintf("%s/yearly-plans/%s/devotionals", ps.apiUrl, planId)

	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := http.Post(endpoint, "json", bytes.NewBuffer(body))
	log.Printf("[%s] %s \n\n", "POST", endpoint)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		log.Printf("ERROR: [%s] %s \npayload: %s\n\n", "POST", endpoint, string(body))
		return ErrAddingDailyDevotional
	}
	return nil
}
