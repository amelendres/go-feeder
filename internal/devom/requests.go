package devom

type AddDailyDevotionalReq struct {
	PlanId       string `json:"-"`
	DevotionalId string `json:"devotionalId"`
	Day          int    `json:"day"`
}

type AddNextDevotionalReq struct {
	PlanId       string `json:"-"`
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
	DevotionalId string `json:"-"`
	TopicId      string `json:"topicId"`
}
