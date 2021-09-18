package devom

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	feed "github.com/amelendres/go-feeder/pkg"
	"github.com/google/uuid"
)

var (
	ErrImportingTopics = errors.New("fails importing topics")
	ErrGettingPlans    = func(want, got int) error {
		return fmt.Errorf("fails getting plans, unexpected response status, want %d but got %d", want, got)
	}
	ErrCreatingTopic = func(want, got int) error {
		return fmt.Errorf("fails creating topic, unexpected response status, want %d but got %d", want, got)
	}
	ErrYearlyPlanNotFound = func(year int) error {
		return fmt.Errorf("Plan <%d> not found", year)
	}
	ErrDailyDevotionalNotFound = func(planId string, day int) error {
		return fmt.Errorf("Daily Devotional not found <%s : %d> not found", planId, day)
	}
)

type TopicSender struct {
	apiUrl      string
	to          Destination
	yearlyPlans map[string]*YearlyPlan
	topics      map[string]*Topic
}

func NewTopicSender(apiUrl string) feed.Sender {
	return &TopicSender{apiUrl: apiUrl}
}

func (ts *TopicSender) Send(items []feed.Feed, to feed.Destination) error {
	ts.to = to.(Destination)
	var errors []error
	for _, item := range items {
		topic := ts.mapItem(item)
		if err := ts.createTopic(*topic); err != nil {
			log.Print(err)
			errors = append(errors, err)
			continue
		}
		err := ts.addTopic(*topic, item["devotionals"])
		if err != nil {
			errors = append(errors, err)
		}
	}

	if errors == nil {
		return nil
	}
	return ErrImportingTopics
}

func (ts *TopicSender) createTopic(topic Topic) error {
	endpoint := fmt.Sprintf("%s/categories", ts.apiUrl)

	body, err := json.Marshal(topic)
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
		err = ErrCreatingTopic(http.StatusCreated, resp.StatusCode)
		logRequestResponseError(req, resp, body, err)
		return err
	}
	logRequest(req)
	return nil
}

func (ts *TopicSender) addTopic(topic Topic, yearlyDevotionalsJSON string) error {

	yealyDevotionals := mapYearlyDevotionalsFromJSON(yearlyDevotionalsJSON)
	var err error
	for _, dev := range yealyDevotionals {
		plan := ts.yearlyPlan(GetYearlyPlanReq{Year: dev.Year, AuthorId: ts.to.AuthorId})
		if plan == nil {
			err = ErrYearlyPlanNotFound(dev.Year)
			// log.Println(err)
			log.Fatal(err)
			continue
			// return ErrYearlyPlanNotFound(year)
		}

		dd := ts.dailyDevotional(GetPlanDevotionalReq{TopicId: plan.TopicId, Day: dev.Day})
		if dd == nil {
			err = ErrDailyDevotionalNotFound(plan.Id, dev.Day)
			log.Println(err)
			continue
		}
		err = ts.addDevotionalTopic(AddDevotionalTopicReq{dd.Devotional.Id, topic.Id})
		if err != nil {
			log.Print(err)
			continue
			// return err
		}
	}
	return err
}

func (ts *TopicSender) addDevotionalTopic(data AddDevotionalTopicReq) error {
	endpoint := fmt.Sprintf("%s/devotionals/%s/topics/add", ts.apiUrl, data.DevotionalId)

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

func (ts *TopicSender) yearlyPlan(getPlan GetYearlyPlanReq) *YearlyPlan {
	//refresh cache
	if ts.yearlyPlans == nil {
		if err := ts.refreshYearlyPlans(getPlan.AuthorId); err != nil {
			log.Print(err)
			return nil
		}
		if err := ts.refreshTopics(); err != nil {
			log.Print(err)
			return nil
		}
	}

	topic, ok := ts.topics[strconv.Itoa(getPlan.Year)]
	if !ok {
		return nil
	}

	if plan, ok := ts.yearlyPlans[topic.Id]; ok {
		return plan
	}

	return nil
}

func (ts *TopicSender) refreshCache(authorId string) error {
	if err := ts.refreshYearlyPlans(authorId); err != nil {
		return err
	}
	if err := ts.refreshTopics(); err != nil {
		return err
	}
	return nil
}

func (ts *TopicSender) refreshYearlyPlans(authorId string) error {
	endpoint := fmt.Sprintf("%s/yearly-plans?authorId=%s", ts.apiUrl, authorId)
	resp, err := get(endpoint)
	if err != nil {
		return err
	}

	yearlyPlans, err := newYearlyPlansFromJSON(resp.Body)
	if err != nil {
		return err
	}

	ts.yearlyPlans = make(map[string]*YearlyPlan)
	for _, plan := range yearlyPlans {
		dailyDevotionals, err := ts.getDailyDevotionals(plan.Id)
		if err != nil {
			log.Println(err)
			continue
		}
		ddIdx := make(map[string]*DailyDevotional)
		for _, dd := range dailyDevotionals {
			ddIdx[string(dd.Day)] = dd
		}
		plan.DailyDevotionals = ddIdx
		ts.yearlyPlans[string(plan.TopicId)] = plan
	}
	return nil
}

func (ts *TopicSender) refreshTopics() error {
	endpoint := fmt.Sprintf("%s/categories", ts.apiUrl)
	resp, err := get(endpoint)
	if err != nil {
		return err
	}
	topics, err := newTopicsFromJSON(resp.Body)
	if err != nil {
		return err
	}

	ts.topics = make(map[string]*Topic)
	for _, topic := range topics {
		ts.topics[topic.Title] = topic
	}
	return nil
}

func (ts *TopicSender) getDailyDevotionals(planId string) ([]*DailyDevotional, error) {
	endpoint := fmt.Sprintf("%s/yearly-plans/%s/devotionals", ts.apiUrl, planId)
	resp, err := get(endpoint)
	if err != nil {
		return nil, err
	}

	dailyDevotionals, err := newDailyDevotionalsFromJSON(resp.Body)
	if err != nil {
		return nil, err
	}

	return dailyDevotionals, nil
}

func newYearlyPlansFromJSON(rdr io.Reader) ([]*YearlyPlan, error) {
	var plans []*YearlyPlan
	err := json.NewDecoder(rdr).Decode(&plans)

	if err != nil {
		err = fmt.Errorf("problem parsing Plans, %+v", err)
	}

	return plans, err
}
func newTopicsFromJSON(rdr io.Reader) ([]*Topic, error) {
	var topics []*Topic
	err := json.NewDecoder(rdr).Decode(&topics)

	if err != nil {
		err = fmt.Errorf("problem parsing Topics, %+v", err)
	}

	return topics, err
}

func newDailyDevotionalsFromJSON(rdr io.Reader) ([]*DailyDevotional, error) {
	var items []*DailyDevotional
	err := json.NewDecoder(rdr).Decode(&items)

	if err != nil {
		err = fmt.Errorf("problem parsing DailyDevotionals, %+v", err)
	}

	return items, err
}

func (ts *TopicSender) dailyDevotional(getDev GetPlanDevotionalReq) *DailyDevotional {
	//from cache
	if dev, ok := ts.yearlyPlans[getDev.TopicId].DailyDevotionals[string(getDev.Day)]; ok {
		return dev
	}
	return nil
}

func (ts *TopicSender) mapItem(feed feed.Feed) *Topic {

	return &Topic{
		uuid.New().String(),
		feed["title"],
		"",
		0,
		ts.to.AuthorId,
	}
}

func mapYearlyDevotionalsFromJSON(txtJSON string) []YearlyDevotional {
	var items []YearlyDevotional
	err := json.Unmarshal([]byte(txtJSON), &items)
	if err != nil {
		log.Fatal(err)
	}
	return items
}

func get(endpoint string) (*http.Response, error) {
	req, err := http.NewRequest("GET", endpoint, bytes.NewBuffer(nil))
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logRequestError(req, err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		err = ErrGettingPlans(http.StatusOK, resp.StatusCode)
		logRequestResponseError(req, resp, nil, err)
		return nil, err
	}
	logRequest(req)
	return resp, nil
}
