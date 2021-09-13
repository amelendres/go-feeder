package devom

type DailyDevotional struct {
	Day        int        `json:"day"`
	Devotional Devotional `json:"devotional"`
}

type YearlyPlan struct {
	Id       string `json:"id"`
	Year     int    `json:"year"`
	AuthorId string `json:"authorId"`
}

type AddDailyDevotionalReq struct {
	Day          int    `json:"day"`
	DevotionalId string `json:"devotionalId"`
}
