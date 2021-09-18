package devom

type AddDailyDevotionalReq struct {
	Day          int    `json:"day"`
	DevotionalId string `json:"devotionalId"`
}

type GetYearlyPlanReq struct {
	Year     int    `json:"year"`
	AuthorId string `json:"authorId"`
}

type GetPlanDevotionalReq struct {
	TopicId string `json:"topicId"`
	Day     int    `json:"day"`
}

type AddDevotionalTopicReq struct {
	DevotionalId string `json:"devotionalId"`
	TopicId      string `json:"topicId"`
}
