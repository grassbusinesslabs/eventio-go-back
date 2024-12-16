package domain

import "time"

type Event struct {
	Id          uint64
	User_Id     uint64
	Title       string
	Description string
	Date        time.Time
	Image       string
	Location    string
	Lat         float64
	Lon         float64
	CreatedDate time.Time
	UpdatedDate time.Time
	DeletedDate *time.Time
}

func (e Event) GetUserId() uint64 {
	return e.User_Id
}
