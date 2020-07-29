package devom

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
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

var ErrCreatingDevotional = errors.New("does not create devotional")

func CreateDevotional(dev Devotional) error {
	url := fmt.Sprintf("%s/devotionals", os.Getenv("DEVOM_API_URL"))
	body, err := json.Marshal(dev)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "json", bytes.NewBuffer(body))
	log.Printf("[%s] %s \n", "POST", url)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("ERROR: [%s] %s \npayload: %s\n\n", "POST", url, string(body))
		return ErrCreatingDevotional
	}

	return nil
}
