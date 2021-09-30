package devom

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	feed "github.com/amelendres/go-feeder/pkg"
	"github.com/google/uuid"
)

var (
	ErrImportingTopics = errors.New("fails importing topics")

	ErrDailyDevotionalNotFound = func(planId string, day int) error {
		return fmt.Errorf("Daily Devotional not found <%s : %d> not found", planId, day)
	}
	ErrYearlyPlanNotFound = func(year int) error {
		return fmt.Errorf("Plan <%d> not found", year)
	}
	ErrTopicPlanNotFound = func(topicId string) error {
		return fmt.Errorf("Plan <%s> not found", topicId)
	}
)

type TopicSender struct {
	api    API
	to     *feed.Destination
	plans  map[string]*Plan
	topics map[string]*Topic
}

func NewTopicSender(api API) feed.Sender {
	return &TopicSender{api: api}
}

func (ts *TopicSender) Destination(d *feed.Destination) {
	ts.to = d
}

func (ts *TopicSender) Send(items []feed.Item) error {
	if err := ts.refreshCache(ts.to.AuthorId); err != nil {
		log.Fatal(err)
		return nil
	}

	var errors []error
	for _, item := range items {
		topic := ts.topic(topicTitle(item))
		if topic == nil {
			topic = ts.mapItem(item)
			//create new topic
			if err := ts.api.createTopic(*topic); err != nil {
				return err
			}
			ts.topics[topic.Title] = topic
		}

		//categorize devotionals
		err := ts.addTopicToDevotionals(*topic, item["devotionals"])
		if err != nil {
			errors = append(errors, err)
			continue
		}

		//create topic plan
		topicPlan := &Plan{
			Id:          uuid.New().String(),
			Title:       planTitle(item),
			Description: "",
			TopicId:     topic.Id,
			AuthorId:    ts.to.AuthorId,
			PublisherId: ts.to.PublisherId,
		}
		err = ts.api.createPlan(*topicPlan)
		if err != nil {
			return err
		}
		ts.plans[topicPlan.TopicId] = topicPlan

		//add devotionals to the topic plan
		err = ts.addDailyDevotionals(*topicPlan, item["devotionals"])
		if err != nil {
			errors = append(errors, err)
			continue
		}
	}

	if errors == nil {
		return nil
	}
	return ErrImportingTopics
}

func (ts *TopicSender) addDailyDevotionals(plan Plan, yearlyDevotionalsJSON string) error {

	yealyDevotionals := ts.mapYearlyDevotionalsFromJSON(yearlyDevotionalsJSON)
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

		err = ts.api.addNextDevotional(AddNextDevotionalReq{PlanId: plan.Id, DevotionalId: dd.Devotional.Id})
		if err != nil {
			log.Println(err)
			continue
		}
	}
	return err
}

func (ts *TopicSender) addTopicToDevotionals(topic Topic, yearlyDevotionalsJSON string) error {

	yealyDevotionals := ts.mapYearlyDevotionalsFromJSON(yearlyDevotionalsJSON)
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
		err = ts.api.addDevotionalTopic(AddDevotionalTopicReq{dd.Devotional.Id, topic.Id})
		if err != nil {
			log.Print(err)
			continue
		}
	}
	return err
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
	plans, err := ts.api.getPlans(authorId)
	if err != nil {
		return err
	}

	ts.plans = make(map[string]*Plan)
	for _, plan := range plans {
		ts.plans[plan.TopicId] = plan
	}
	return nil
}

func (ts *TopicSender) refreshTopics() error {
	topics, err := ts.api.getTopics()
	if err != nil {
		return err
	}

	ts.topics = make(map[string]*Topic)
	for _, topic := range topics {
		ts.topics[topic.Title] = topic
	}
	return nil
}

func (ts *TopicSender) topic(title string) *Topic {
	//from cache
	if topic, ok := ts.topics[title]; ok {
		return topic
	}
	return nil
}

func (ts *TopicSender) dailyDevotional(getDev GetPlanDevotionalReq) *DailyDevotional {
	//from cache
	if dev, ok := ts.plans[getDev.TopicId].DailyDevotionals[fmt.Sprint(getDev.Day)]; ok {
		return dev
	}
	return nil
}

func (ts *TopicSender) mapItem(feed feed.Item) *Topic {

	return &Topic{
		uuid.New().String(),
		topicTitle(feed),
		"",
		0,
		ts.to.AuthorId,
	}
}

func topicTitle(feed feed.Item) string {
	return strings.Split(feed["title"], ",")[0]
}

func planTitle(feed feed.Item) string {
	txt := strings.Split(feed["title"], ",")
	if len(txt) > 1 {
		return txt[1] + " " + txt[0]
	}
	return txt[0]
}

func (ts *TopicSender) mapYearlyDevotionalsFromJSON(txtJSON string) []YearlyDevotional {
	var items []YearlyDevotional
	err := json.Unmarshal([]byte(txtJSON), &items)
	if err != nil {
		log.Fatal(err)
	}
	return items
}
