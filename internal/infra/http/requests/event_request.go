package requests

import (
	"time"

	"github.com/grassbusinesslabs/eventio-go-back/internal/domain"
)

type EventRequest struct {
	Tytle       string    `json:"tytle" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Date        time.Time `json:"date" validate:"required"`
	Image       string    `json:"image"`
	Location    string    `json:"location" validate:"required"`
	Lat         float64   `json:"lat" validate:"required"`
	Lon         float64   `json:"lon" validate:"required"`
}

func (r EventRequest) ToDomainModel() (interface{}, error) {
	return domain.Event{
		Title:       r.Tytle,
		Description: r.Description,
		Date:        r.Date,
		Image:       r.Image,
		Location:    r.Location,
		Lat:         r.Lat,
		Lon:         r.Lon,
	}, nil
}
