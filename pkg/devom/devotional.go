package devom

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


func NewPassage(text, reference string) Passage {
	return Passage{text, reference}
}

//func sendDevotional(dev Devotional) error {
//	url := fmt.Sprintf("%s/devotionals", os.Getenv("DEVOM_API_URL"))
//	body, err := json.Marshal(dev)
//	if err != nil {
//		return err
//	}
//
//	resp, err := http.Post(url, "json", bytes.NewBuffer(body))
//	log.Printf("[%s] %s \n", "POST", url)
//
//	if err != nil {
//		log.Printf("ERROR: [%s] %s \npayload: %s\n\n", "POST", url, string(body))
//		return err
//	}
//
//	if resp.StatusCode != http.StatusOK {
//		log.Printf("STATUS ERROR: [%s] %s \npayload: %s\n\n reponse: status %d", "POST", url, string(body), resp.StatusCode)
//		return ErrSendingDevotional
//	}
//
//	return nil
//}
