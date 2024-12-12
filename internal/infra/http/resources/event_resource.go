package resources

import (
	"time"

	"github.com/grassbusinesslabs/eventio-go-back/internal/domain"
)

type EventDto struct {
	EventId     uint64    `json:"eventid"`
	Tytle       string    `json:"tytle"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Image       string    `json:"image"`
	Location    string    `json:"location"`
	Lat         float64   `json:"lat"`
	Lon         float64   `json:"lon"`
}

type ShortEventDto struct {
	EventId  uint64    `json:"eventid"`
	Tytle    string    `json:"tytle"`
	Date     time.Time `json:"date"`
	Location string    `json:"location"`
	Lat      float64   `json:"lat"`
	Lon      float64   `json:"lon"`
}

type EventsDto struct {
	Events []ShortEventDto `json:"events"`
}

func (d EventsDto) DomainToDto(ev []domain.Event) EventsDto {
	events := make([]ShortEventDto, len(ev))
	for i, e := range ev {
		events[i] = ShortEventDto{}.DomainToDto(e)
	}

	return EventsDto{
		Events: events,
	}
}

func (d EventDto) DomainToDto(t domain.Event) EventDto {
	return EventDto{
		EventId:     t.EventId,
		Tytle:       t.Title,
		Description: t.Description,
		Date:        t.Date,
		Image:       t.Image,
		Location:    t.Location,
		Lat:         t.Lat,
		Lon:         t.Lon,
	}
}

func (d ShortEventDto) DomainToDto(t domain.Event) ShortEventDto {
	return ShortEventDto{
		EventId:  t.EventId,
		Tytle:    t.Title,
		Date:     t.Date,
		Location: t.Location,
		Lat:      t.Lat,
		Lon:      t.Lon,
	}
}
