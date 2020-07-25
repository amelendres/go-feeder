package devom

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type YearlyPlan struct {
	Id   string `json:"id"`
	Year int    `json:"year"`
}

type DailyDevotional struct {
	Day          int    `json:"day"`
	DevotionalId string `json:"devotionalId"`
}

var ErrAddingDailyDevotional = errors.New("does not add daily devotional")

func AddDailyDevotional(dev DailyDevotional, planId string) error {
	url := fmt.Sprintf("http://localhost:8030/api/v1/yearly-plans/%s/devotionals", planId)
	body, err := json.Marshal(dev)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "json", bytes.NewBuffer(body))
	log.Printf("[%s] %s \n\n", "POST", url)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("ERROR: [%s] %s \npayload: %s\n\n", "POST", url, string(body))
		return ErrAddingDailyDevotional
	}
	return nil
}
