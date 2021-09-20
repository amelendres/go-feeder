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
	"strings"

	feed "github.com/amelendres/go-feeder/pkg"
	"github.com/google/uuid"
)

var (
	ErrImportingTopics = errors.New("fails importing topics")
	ErrGettingResource = func(want, got int) error {
		return fmt.Errorf("fails getting resource, unexpected response status, want %d but got %d", want, got)
	}
	ErrCreatingResource = func(want, got int) error {
		return fmt.Errorf("fails creating resource, unexpected response status, want %d but got %d", want, got)
	}
	ErrYearlyPlanNotFound = func(year int) error {
		return fmt.Errorf("Plan <%d> not found", year)
	}
	ErrTopicPlanNotFound = func(topicId string) error {
		return fmt.Errorf("Plan <%s> not found", topicId)
	}
	ErrDailyDevotionalNotFound = func(planId string, day int) error {
		return fmt.Errorf("Daily Devotional not found <%s : %d> not found", planId, day)
	}
)

type TopicSender struct {
	apiUrl string
	to     Destination
	plans  map[string]*Plan
	topics map[string]*Topic
}

func NewTopicSender(apiUrl string) feed.Sender {
	return &TopicSender{apiUrl: apiUrl}
}

func (ts *TopicSender) Send(items []feed.Feed, to feed.Destination) error {
	ts.to = to.(Destination)
	if err := ts.refreshCache(ts.to.AuthorId); err != nil {
		log.Fatal(err)
		return nil
	}

	var errors []error
	for _, item := range items {
		topic := ts.mapItem(item)
		if err := ts.createTopic(*topic); err != nil {
			// log.Println(err)
			log.Fatal(err)
			errors = append(errors, err)
			continue
		}
		err := ts.addTopicToDevotionals(*topic, item["devotionals"])
		if err != nil {
			log.Fatal(err)
			errors = append(errors, err)
			continue
		}

		topicPlan, err := ts.createTopicPlan(*topic)
		if err != nil {
			log.Fatal(err)
			// log.Println(err)
			errors = append(errors, err)
			continue
		}
		err = ts.addDailyDevotionals(*topicPlan, item["devotionals"])
		if err != nil {
			log.Fatal(err)
			errors = append(errors, err)
			continue
		}
	}

	if errors == nil {
		return nil
	}
	return ErrImportingTopics
}

func (ts *TopicSender) createTopic(topic Topic) error {
	endpoint := fmt.Sprintf("%s/categories", ts.apiUrl)
	_, err := post(endpoint, topic)
	if err != nil {
		return err
	}

	ts.topics[topic.Title] = &topic

	return nil
}

func (ts *TopicSender) createTopicPlan(topic Topic) (*Plan, error) {
	endpoint := fmt.Sprintf("%s/yearly-plans", ts.apiUrl)
	plan := &Plan{
		Id:          uuid.New().String(),
		Title:       topic.Title,
		Description: "Plan importado",
		TopicId:     topic.Id,
		AuthorId:    ts.to.AuthorId,
		PublisherId: ts.to.PublisherId,
	}
	_, err := post(endpoint, plan)
	if err != nil {
		return nil, err
	}
	ts.plans[plan.TopicId] = plan

	return plan, nil
}

func (ts *TopicSender) addDailyDevotionals(plan Plan, yearlyDevotionalsJSON string) error {

	yealyDevotionals := mapYearlyDevotionalsFromJSON(yearlyDevotionalsJSON)
	var err error
	for _, dev := range yealyDevotionals {
		yearlyPlan := ts.yearlyPlan(GetYearlyPlanReq{Year: dev.Year, AuthorId: ts.to.AuthorId})
		if yearlyPlan == nil {
			err = ErrYearlyPlanNotFound(dev.Year)
			log.Println(err)
			continue
		}

		dd := ts.dailyDevotional(GetPlanDevotionalReq{TopicId: yearlyPlan.TopicId, Day: dev.Day})
		if dd == nil {
			err = ErrDailyDevotionalNotFound(plan.Id, dev.Day)
			log.Println(err)
			continue
		}

		err = ts.addNextDevotional(AddNextDevotionalReq{PlanId: plan.Id, DevotionalId: dd.Devotional.Id})
		if err != nil {
			log.Println(err)
			continue
		}
	}
	return err
}

func (ts *TopicSender) addTopicToDevotionals(topic Topic, yearlyDevotionalsJSON string) error {

	yealyDevotionals := mapYearlyDevotionalsFromJSON(yearlyDevotionalsJSON)
	var err error
	for _, dev := range yealyDevotionals {
		plan := ts.yearlyPlan(GetYearlyPlanReq{Year: dev.Year, AuthorId: ts.to.AuthorId})
		if plan == nil {
			err = ErrYearlyPlanNotFound(dev.Year)
			log.Println(err)
			continue
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
		}
	}
	return err
}

func (ts *TopicSender) addDevotionalTopic(body AddDevotionalTopicReq) error {
	endpoint := fmt.Sprintf("%s/devotionals/%s/topics/add", ts.apiUrl, body.DevotionalId)
	_, err := post(endpoint, body)
	if err != nil {
		return err
	}

	return nil
}

func (ts *TopicSender) addNextDevotional(body AddNextDevotionalReq) error {

	endpoint := fmt.Sprintf("%s/yearly-plans/%s/devotionals", ts.apiUrl, body.PlanId)
	_, err := post(endpoint, body)
	if err != nil {
		return err
	}

	return nil
}

func (ts *TopicSender) yearlyPlan(getPlan GetYearlyPlanReq) *Plan {

	topic, ok := ts.topics[strconv.Itoa(getPlan.Year)]
	if !ok {
		return nil
	}

	if plan, ok := ts.plans[topic.Id]; ok {
		return plan
	}

	return nil
}

func (ts *TopicSender) refreshCache(authorId string) error {
	if err := ts.refreshPlans(authorId); err != nil {
		return err
	}
	if err := ts.refreshTopics(); err != nil {
		return err
	}
	return nil
}

func (ts *TopicSender) refreshPlans(authorId string) error {
	endpoint := fmt.Sprintf("%s/yearly-plans?authorId=%s", ts.apiUrl, authorId)
	resp, err := get(endpoint)
	if err != nil {
		return err
	}

	yearlyPlans, err := newPlansFromJSON(resp.Body)
	if err != nil {
		return err
	}

	ts.plans = make(map[string]*Plan)
	for _, plan := range yearlyPlans {
		dailyDevotionals, err := ts.getDailyDevotionals(plan.Id)
		if err != nil {
			log.Println(err)
			continue
		}
		ddIdx := make(map[string]*DailyDevotional)
		for _, dd := range dailyDevotionals {
			ddIdx[fmt.Sprint(dd.Day)] = dd
		}
		plan.DailyDevotionals = ddIdx
		ts.plans[plan.TopicId] = plan
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

func newPlansFromJSON(rdr io.Reader) ([]*Plan, error) {
	var plans []*Plan
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
	if dev, ok := ts.plans[getDev.TopicId].DailyDevotionals[fmt.Sprint(getDev.Day)]; ok {
		return dev
	}
	return nil
}

func (ts *TopicSender) mapItem(feed feed.Feed) *Topic {

	return &Topic{
		uuid.New().String(),
		strings.Split(feed["title"], ",")[0],
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
		err = ErrGettingResource(http.StatusOK, resp.StatusCode)
		logRequestResponseError(req, resp, nil, err)
		return nil, err
	}
	logRequest(req)
	return resp, nil
}

func post(endpoint string, obj interface{}) (*http.Response, error) {
	body, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logRequestError(req, err)
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		err = ErrCreatingResource(http.StatusCreated, resp.StatusCode)
		logRequestResponseError(req, resp, body, err)
		return nil, err
	}
	logRequest(req)
	return resp, nil
}
