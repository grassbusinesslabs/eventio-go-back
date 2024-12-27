package resources

import (
	"time"

	"github.com/grassbusinesslabs/eventio-go-back/internal/domain"
)

type EventDto struct {
	Id          uint64    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Image       string    `json:"image"`
	City        string    `json:"city"`
	Location    string    `json:"location"`
	Lat         float64   `json:"lat"`
	Lon         float64   `json:"lon"`
	Count       uint64    `json:"count,omitempty"`
}

type ShortEventDto struct {
	Id          uint64    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Image       string    `json:"image"`
	City        string    `json:"city"`
	Location    string    `json:"location"`
	Lat         float64   `json:"lat"`
	Lon         float64   `json:"lon"`
	Count       uint64    `json:"count"`
}

type EventsDto struct {
	Events []ShortEventDto `json:"events"`
}

func (d EventsDto) DomainToDto(ev []domain.Event, count []uint64) EventsDto {
	events := make([]ShortEventDto, len(ev))
	for i, e := range ev {
		events[i] = ShortEventDto{}.DomainToDto(e, count[i])
	}

	return EventsDto{
		Events: events,
	}
}

func (d EventDto) DomainToDto(t domain.Event) EventDto {
	return EventDto{
		Id:          t.Id,
		Title:       t.Title,
		Description: t.Description,
		Date:        t.Date,
		Image:       t.Image,
		City:        t.City,
		Location:    t.Location,
		Lat:         t.Lat,
		Lon:         t.Lon,
	}
}

func (d EventDto) DomainToDtoWC(t domain.Event, count uint64) EventDto {
	return EventDto{
		Id:          t.Id,
		Title:       t.Title,
		Description: t.Description,
		Date:        t.Date,
		Image:       t.Image,
		City:        t.City,
		Location:    t.Location,
		Lat:         t.Lat,
		Lon:         t.Lon,
		Count:       count,
	}
}

func (d ShortEventDto) DomainToDto(t domain.Event, count uint64) ShortEventDto {
	return ShortEventDto{
		Id:          t.Id,
		Title:       t.Title,
		Description: t.Description,
		Date:        t.Date,
		Image:       t.Image,
		City:        t.City,
		Location:    t.Location,
		Lat:         t.Lat,
		Lon:         t.Lon,
		Count:       count,
	}
}
