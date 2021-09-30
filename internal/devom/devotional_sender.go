package devom

import (
	"fmt"
	"log"
	"strconv"

	feed "github.com/amelendres/go-feeder/pkg"
	"github.com/google/uuid"
)

var (
	ErrAddingDailyDevotional = func(want, got int) error {
		return fmt.Errorf("fails adding daily devotional, unexpected response status, want %d but got %d", want, got)
	}
	ErrCreatingDevotional = func(want, got int) error {
		return fmt.Errorf("fails creating devotional, unexpected response status, want %d but got %d", want, got)
	}
)

type devotionalSender struct {
	api         API
	to          *feed.Destination
	plan        *Plan
	devotionals map[string]*Devotional
}

func NewDevotionalSender(api API) feed.Sender {
	return &devotionalSender{api: api}
}

func (ps *devotionalSender) Destination(d *feed.Destination) {
	ps.to = d
}

func (ps *devotionalSender) Send(feeds []feed.Item) error {
	if err := ps.refreshCache(); err != nil {
		log.Fatal(err)
		return nil
	}

	var err error
	for _, f := range feeds {
		dev := ps.mapItem(f)
		day, _ := strconv.Atoi(f["day"])

		//adding current Devotional as Daily Devotional
		if currentDev := ps.devotional(dev.Title); currentDev != nil {
			if dd := ps.dailyDevotional(currentDev.Id); dd != nil {
				continue
			}
			_ = ps.api.addDailyDevotional(AddDailyDevotionalReq{ps.to.PlanId, currentDev.Id, day})
			continue
		}

		if err = ps.api.createDevotional(dev); err != nil {
			return err
		}
		err = ps.api.addDailyDevotional(AddDailyDevotionalReq{ps.to.PlanId, dev.Id, day})
		if err != nil {
			return err
		}
	}
	return err
}

func (ps *devotionalSender) mapItem(feed feed.Item) Devotional {

	return Devotional{
		uuid.New().String(),
		feed["title"],
		Passage{feed["passage_text"], feed["passage_reference"]},
		feed["content"],
		feed["bible_reading"],
		nil,
		ps.to.AuthorId,
		ps.to.PublisherId,
		nil,
	}
}

func (ps *devotionalSender) refreshCache() error {
	plan, err := ps.api.getPlan(ps.to.PlanId)
	if err != nil {
		return err
	}
	ps.plan = plan

	devotionals, err := ps.api.getDevotionals(ps.to.AuthorId)
	if err != nil {
		return err
	}
	ps.devotionals = make(map[string]*Devotional)
	for _, dev := range devotionals {
		ps.devotionals[dev.Title] = dev
	}

	return nil
}

func (ps *devotionalSender) dailyDevotional(id string) *DailyDevotional {
	//from cache
	if dd, ok := ps.plan.DailyDevotionals[id]; ok {
		return dd
	}
	return nil
}

func (ps *devotionalSender) devotional(title string) *Devotional {
	//from cache
	if dev, ok := ps.devotionals[title]; ok {
		return dev
	}
	return nil
}
