package devom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	feed "github.com/amelendres/go-feeder/pkg"
	"github.com/google/uuid"
)

var (
	ErrAddingDailyDevotional = func(want, got int) error {
		return fmt.Errorf("fails adding daily devotional, unexpected response status, want %d but got %d", want, got)
	}
	ErrCreatingDevotional = func(want, got int) error {
		return fmt.Errorf("fails creating devotional, unexpected response status, want %d but got %d", want, got)
	}
)

// TODO: rename to DevotionalSender
type PlanSender struct {
	to     Destination
	apiUrl string
}

func NewPlanSender(apiUrl string) feed.Sender {
	return &PlanSender{apiUrl: apiUrl}
}

func (ps *PlanSender) Send(feeds []feed.Feed, to feed.Destination) error {
	ps.to = to.(Destination)
	var err error
	for _, f := range feeds {
		dev := ps.mapFeed(f)
		if err = ps.sendDevotional(dev); err == nil {
			day, _ := strconv.Atoi(f["day"])
			err = ps.addDailyDevotional(AddDailyDevotionalReq{day, dev.Id}, ps.to.PlanId)
			if err != nil {
				return err
			}
		}
	}
	return err
}

func (ps *PlanSender) mapFeed(feed feed.Feed) Devotional {

	return Devotional{
		uuid.New().String(),
		feed["title"],
		Passage{feed["passage_text"], feed["passage_reference"]},
		feed["content"],
		feed["bible_reading"],
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

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logRequestError(req, err)
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		err = ErrCreatingDevotional(http.StatusCreated, resp.StatusCode)
		logRequestResponseError(req, resp, body, err)
		return err
	}
	logRequest(req)
	return nil
}

func (ps *PlanSender) addDailyDevotional(data AddDailyDevotionalReq, planId string) error {
	endpoint := fmt.Sprintf("%s/yearly-plans/%s/devotionals", ps.apiUrl, planId)

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logRequestError(req, err)
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		err = ErrAddingDailyDevotional(http.StatusCreated, resp.StatusCode)
		logRequestResponseError(req, resp, body, err)
		return err
	}
	logRequest(req)
	return nil
}

func logRequestError(req *http.Request, err error) {
	log.Printf("%s\nrequest: [%s] %s \n", err.Error(), req.Method, req.URL)
}

func logRequest(req *http.Request) {
	log.Printf("[%s] %s \n", req.Method, req.URL)
}

func logRequestResponseError(req *http.Request, resp *http.Response, body []byte, err error) {
	log.Printf(
		"[%s] 😱 %s \n%s\npayload: %s\nreponse status: %d",
		req.Method,
		req.URL,
		err.Error(),
		string(body),
		resp.StatusCode)
}
