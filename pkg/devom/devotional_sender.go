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

// TODO: rename to DevotionalSender
type PlanSender struct {
	api         API
	to          Destination
	plan        *Plan
	devotionals map[string]*Devotional
}

func NewPlanSender(api API) feed.Sender {
	return &PlanSender{api: api}
}

func (ps *PlanSender) Send(feeds []feed.Feed, to feed.Destination) error {
	ps.to = to.(Destination)

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

func (ps *PlanSender) mapItem(feed feed.Feed) Devotional {

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

func (ps *PlanSender) refreshCache() error {
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

func (ps *PlanSender) dailyDevotional(id string) *DailyDevotional {
	//from cache
	if dd, ok := ps.plan.DailyDevotionals[id]; ok {
		return dd
	}
	return nil
}

func (ps *PlanSender) devotional(title string) *Devotional {
	//from cache
	if dev, ok := ps.devotionals[title]; ok {
		return dev
	}
	return nil
}
