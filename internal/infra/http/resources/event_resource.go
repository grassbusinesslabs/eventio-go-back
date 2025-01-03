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
	Pages  uint64          `json:"pages"`
	Total  uint64          `json:"total"`
}

type EventsDtoWC struct {
	Events []ShortEventDtoWC `json:"events"`
	Pages  uint64            `json:"pages"`
	Total  uint64            `json:"total"`
}

func (d EventsDto) DomainToDto(ev domain.Events) EventsDto {
	events := make([]ShortEventDto, len(ev.Items))
	for i, e := range ev.Items {
		events[i] = ShortEventDto{}.DomainToDto(e)
	}

	return EventsDto{
		Events: events,
		Pages:  ev.Pages,
		Total:  ev.Total,
	}
}

func (d EventsDtoWC) DomainToDtoWC(ev domain.Events, count []uint64) EventsDtoWC {
	events := make([]ShortEventDtoWC, len(ev.Items))
	for i, e := range ev.Items {
		events[i] = ShortEventDtoWC{}.DomainToDtoWC(e, count[i])
	}

	return EventsDtoWC{
		Events: events,
		Pages:  ev.Pages,
		Total:  ev.Total,
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
