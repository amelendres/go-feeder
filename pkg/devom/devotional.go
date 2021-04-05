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
