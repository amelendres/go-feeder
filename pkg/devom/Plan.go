package devom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func AddDailyDevotional(dev DailyDevotional, planId string) bool {
	url := fmt.Sprintf("http://localhost:8030/api/v1/yearly-plans/%s/devotionals", planId)
	body, err := json.Marshal(dev)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post(url, "json", bytes.NewBuffer(body))

	if err != nil {
		panic(err)
	}
	//defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("-> %s %s \npayload: %s\n", "POST", url, string(body))

		body, _ = ioutil.ReadAll(resp.Body)
		panic(body)
		return false
	}
	return true
}
