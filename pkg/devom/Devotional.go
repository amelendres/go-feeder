package devom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Devotional struct {
	Id           string   `json:"id"`
	Title        string   `json:"title"`
	Passage      Passage  `json:"passage"`
	Content      string   `json:"content"`
	BibleReading string   `json:"bibleReading"`
	AudioUrl     *string  `json:"audioUrl"`
	AuthorId     string   `json:"authorId"`
	PublisherId  string   `json:"publisherId"`
	Topics       []string `json:"topics"`
}

type Passage struct {
	Text      string `json:"text"`
	Reference string `json:"reference"`
}

func CreateDevotional(dev Devotional) bool {
	url := "http://localhost:8030/api/v1/devotionals"
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
