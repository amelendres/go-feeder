package devom

type DailyDevotional struct {
	Day        int        `json:"day"`
	Devotional Devotional `json:"devotional"`
}

type YearlyPlan struct {
	Id               string `json:"id"`
	TopicId          string `json:"topicId"`
	AuthorId         string `json:"authorId"`
	DailyDevotionals map[string]*DailyDevotional
}
