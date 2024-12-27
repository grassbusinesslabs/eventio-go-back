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
}

type ShortEventDtoWC struct {
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

type EventsDtoWC struct {
	Events []ShortEventDtoWC `json:"events"`
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

func (d EventsDtoWC) DomainToDtoWC(ev []domain.Event, count []uint64) EventsDtoWC {
	events := make([]ShortEventDtoWC, len(ev))
	for i, e := range ev {
		events[i] = ShortEventDtoWC{}.DomainToDtoWC(e, count[i])
	}

	return EventsDtoWC{
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

func (d ShortEventDto) DomainToDto(t domain.Event) ShortEventDto {
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
	}
}

func (d ShortEventDtoWC) DomainToDtoWC(t domain.Event, count uint64) ShortEventDtoWC {
	return ShortEventDtoWC{
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
