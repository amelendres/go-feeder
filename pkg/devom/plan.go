package devom

//type YearlyPlan struct {
//	Id   string `json:"id"`
//	PublisherId string `json:"publisher_id"`
//	AuthorId string `json:"author_id"`
//}
//
//func NewYearlyPlan(id, publisherId, authorId string) feed.Destination  {
//	return &YearlyPlan{id, publisherId, authorId}
//}

type DailyDevotional struct {
	Day          int    `json:"day"`
	DevotionalId string `json:"devotionalId"`
}


//func addDailyDevotional(dev DailyDevotional, PlanId string) error {
//	url := fmt.Sprintf("%s/yearly-plans/%s/devotionals", os.Getenv("DEVOM_API_URL"), PlanId)
//	body, err := json.Marshal(dev)
//	if err != nil {
//		return err
//	}
//
//	resp, err := http.Post(url, "json", bytes.NewBuffer(body))
//	log.Printf("[%s] %s \n\n", "POST", url)
//
//	if err != nil {
//		return err
//	}
//
//	if resp.StatusCode != http.StatusOK {
//		log.Printf("ERROR: [%s] %s \npayload: %s\n\n", "POST", url, string(body))
//		return ErrAddingDailyDevotional
//	}
//	return nil
//}
