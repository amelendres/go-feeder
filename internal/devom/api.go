package devom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

var (
	ErrGettingResource = func(want, got int) error {
		return fmt.Errorf("fails getting resource, unexpected response status, want %d but got %d", want, got)
	}
	ErrCreatingResource = func(want, got int) error {
		return fmt.Errorf("fails creating resource, unexpected response status, want %d but got %d", want, got)
	}
)

type API struct {
	apiUrl string
}

func NewAPI(apiUrl string) *API {
	return &API{apiUrl: apiUrl}
}

// Creates Devotional
func (a *API) createDevotional(dev Devotional) error {
	endpoint := fmt.Sprintf("%s/devotionals", a.apiUrl)
	_, err := post(endpoint, dev)
	if err != nil {
		return err
	}

	return nil
}

func (a *API) getDevotionals(authorId string) ([]*Devotional, error) {
	endpoint := fmt.Sprintf("%s/devotionals?authorId=%s", a.apiUrl, authorId)
	resp, err := get(endpoint)
	if err != nil {
		return nil, err
	}
	items, err := newDevotionalsFromJSON(resp.Body)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (a *API) addDevotionalTopic(req AddDevotionalTopicReq) error {
	endpoint := fmt.Sprintf("%s/devotionals/%s/topics/add", a.apiUrl, req.DevotionalId)
	_, err := post(endpoint, req)
	if err != nil {
		return err
	}

	return nil
}

// Creates Plan
func (a *API) createPlan(plan Plan) error {
	endpoint := fmt.Sprintf("%s/yearly-plans", a.apiUrl)
	_, err := post(endpoint, plan)
	if err != nil {
		return err
	}

	return nil
}

func (a *API) addDailyDevotional(req AddDailyDevotionalReq) error {
	endpoint := fmt.Sprintf("%s/yearly-plans/%s/devotionals", a.apiUrl, req.PlanId)
	_, err := post(endpoint, req)
	if err != nil {
		return err
	}

	return nil
}

func (a *API) addNextDevotional(body AddNextDevotionalReq) error {
	endpoint := fmt.Sprintf("%s/yearly-plans/%s/devotionals", a.apiUrl, body.PlanId)
	_, err := post(endpoint, body)
	if err != nil {
		return err
	}

	return nil
}

func (a *API) getPlans(authorId string) ([]*Plan, error) {
	endpoint := fmt.Sprintf("%s/yearly-plans?authorId=%s", a.apiUrl, authorId)
	resp, err := get(endpoint)
	if err != nil {
		return nil, err
	}

	plans, err := newPlansFromJSON(resp.Body)
	if err != nil {
		return nil, err
	}
	for _, item := range plans {
		dailyDevotionals, err := a.getDailyDevotionals(item.Id)
		if err != nil {
			log.Println(err)
			continue
		}
		ddIdx := make(map[string]*DailyDevotional)
		for _, dd := range dailyDevotionals {
			ddIdx[fmt.Sprint(dd.Day)] = dd
		}
		item.DailyDevotionals = ddIdx
	}
	return plans, nil
}

func (a *API) getPlan(planId string) (*Plan, error) {
	endpoint := fmt.Sprintf("%s/yearly-plans/%s", a.apiUrl, planId)
	resp, err := get(endpoint)
	if err != nil {
		return nil, err
	}

	plan, err := newPlanFromJSON(resp.Body)
	if err != nil {
		return nil, err
	}

	dailyDevotionals, err := a.getDailyDevotionals(plan.Id)
	if err != nil {
		return nil, err
	}
	ddIdx := make(map[string]*DailyDevotional)
	for _, dd := range dailyDevotionals {
		ddIdx[dd.Devotional.Id] = dd
	}
	plan.DailyDevotionals = ddIdx

	return plan, nil
}

func (a *API) getDailyDevotionals(planId string) ([]*DailyDevotional, error) {
	endpoint := fmt.Sprintf("%s/yearly-plans/%s/devotionals", a.apiUrl, planId)
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

// Creates Topic
func (a *API) createTopic(topic Topic) error {
	endpoint := fmt.Sprintf("%s/categories", a.apiUrl)
	_, err := post(endpoint, topic)
	if err != nil {
		return err
	}

	return nil
}

func (a *API) getTopics() ([]*Topic, error) {
	endpoint := fmt.Sprintf("%s/categories", a.apiUrl)
	resp, err := get(endpoint)
	if err != nil {
		return nil, err
	}
	topics, err := newTopicsFromJSON(resp.Body)
	if err != nil {
		return nil, err
	}
	return topics, nil
}

func newDevotionalsFromJSON(rdr io.Reader) ([]*Devotional, error) {
	var devotionals []*Devotional
	err := json.NewDecoder(rdr).Decode(&devotionals)

	if err != nil {
		err = fmt.Errorf("problem parsing Devotionals, %+v", err)
	}

	return devotionals, err
}

func newPlanFromJSON(rdr io.Reader) (*Plan, error) {
	var plan *Plan
	err := json.NewDecoder(rdr).Decode(&plan)

	if err != nil {
		err = fmt.Errorf("problem parsing Plan, %+v", err)
	}

	return plan, err
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

func logRequestError(req *http.Request, err error) {
	log.Printf("[%s] 😱 %s\nerror:%s\n", req.Method, req.URL, err.Error())
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
