package devom

type DailyDevotional struct {
	Day        int        `json:"day"`
	Devotional Devotional `json:"devotional"`
}

type Plan struct {
	Id               string                      `json:"id"`
	Title            string                      `json:"title"`
	Description      string                      `json:"description"`
	CoverPhotoUrl    string                      `json:"coverPhotoUrl"`
	TopicId          string                      `json:"topicId"`
	AuthorId         string                      `json:"authorId"`
	PublisherId      string                      `json:"publisherId"`
	DailyDevotionals map[string]*DailyDevotional `json:"-"`
}
