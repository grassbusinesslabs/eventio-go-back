package requests

import (
	"time"

	"github.com/grassbusinesslabs/eventio-go-back/internal/domain"
)

type EventRequest struct {
	Title       string  `json:"title" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Date        int64   `json:"date" validate:"required"`
	Image       string  `json:"image"`
	Location    string  `json:"location" validate:"required"`
	Lat         float64 `json:"lat" validate:"required"`
	Lon         float64 `json:"lon" validate:"required"`
}

func (r EventRequest) ToDomainModel() (interface{}, error) {
	date := time.Unix(r.Date, 0)
	return domain.Event{
		Title:       r.Title,
		Description: r.Description,
		Date:        date,
		Image:       r.Image,
		Location:    r.Location,
		Lat:         r.Lat,
		Lon:         r.Lon,
	}, nil
}
