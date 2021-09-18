package devom

type Topic struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Position    int    `json:"position"`
	AuthorId    string `json:"authorId"`
}

type YearlyDevotional struct {
	Year int
	Day  int
}
