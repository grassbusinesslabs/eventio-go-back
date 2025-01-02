package database

import (
	"strings"
	"time"

	"github.com/grassbusinesslabs/eventio-go-back/internal/domain"
	"github.com/upper/db/v4"
)

const EventsTableName = "events"

type event struct {
	Id          uint64     `db:"id,omitempty"`
	User_Id     uint64     `db:"user_id"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	Date        time.Time  `db:"date"`
	Image       string     `db:"image"`
	City        string     `db:"city"`
	Location    string     `db:"location"`
	Lat         float64    `db:"lat"`
	Lon         float64    `db:"lon"`
	CreatedDate time.Time  `db:"created_date"`
	UpdatedDate time.Time  `db:"updated_date"`
	DeletedDate *time.Time `db:"deleted_date,omitempty"`
}

type EventSearchParams struct {
	User_Id    uint64
	City       string
	DateDay    *time.Time
	DateMonth  *time.Time
	DateYear   *time.Time
	Search     string
	Location   string
	Pagination domain.Pagination
}

type EventRepository struct {
	coll db.Collection
	sess db.Session
}

func NewEventRepository(sess db.Session) EventRepository {
	return EventRepository{
		coll: sess.Collection(EventsTableName),
		sess: sess,
	}
}

func (r EventRepository) Save(t domain.Event) (domain.Event, error) {
	evn := r.mapDomainToModel(t)
	evn.CreatedDate = time.Now()
	evn.UpdatedDate = time.Now()

	err := r.coll.InsertReturning(&evn)
	if err != nil {
		return domain.Event{}, err
	}

	t = r.mapModelToDomain(evn)
	return t, nil
}

func (r EventRepository) FindById(id uint64) (domain.Event, error) {
	var evn event
	err := r.coll.Find(db.Cond{"id": id}).One(&evn)
	if err != nil {
		return domain.Event{}, err
	}
	ev := r.mapModelToDomain(evn)
	return ev, nil
}

func (r EventRepository) FindListBy(str EventSearchParams) (domain.Events, error) {
	query := r.coll.Find(db.Cond{"deleted_date": nil})

	if str.User_Id != 0 {
		query = query.And(db.Cond{"user_id": str.User_Id})
	}

	if str.City != "" {
		query = query.And(db.Cond{"city": str.City})
	}

	if str.DateDay != nil {
		startDay := time.Date(str.DateDay.Year(), str.DateDay.Month(), str.DateDay.Day(), 0, 0, 0, 0, str.DateDay.Location())
		endDay := startDay.Add(24 * time.Hour)
		query = query.And(db.Cond{"date >=": startDay, "date <": endDay})
	}

	if str.DateMonth != nil {
		startMonth := time.Date(str.DateMonth.Year(), str.DateMonth.Month(), 0, 0, 0, 0, 0, str.DateMonth.Location())
		endMonth := startMonth.AddDate(0, 1, 0)
		query = query.And(db.Cond{"date >=": startMonth, "date <": endMonth})
	}

	if str.DateYear != nil {
		startYear := time.Date(str.DateYear.Year(), 1, 0, 0, 0, 0, 0, str.DateYear.Location())
		endYear := startYear.AddDate(1, 0, 0)
		query = query.And(db.Cond{"date >=": startYear, "date <": endYear})
	}

	if str.Search != "" {
		search := "%" + strings.ToLower(str.Search) + "%"
		query = query.And(db.Raw(`(LOWER(title) LIKE ? OR LOWER(description) LIKE ?)`, search, search))
	}

	if str.Location != "" {
		location := "%" + strings.ToLower(str.Location) + "%"
		query = query.And(db.Raw(`(LOWER(location) LIKE ?)`, location))
	}

	paginate := query.Paginate(uint(str.Pagination.CountPerPage))
	var evns []event

	err := paginate.Page(uint(str.Pagination.Page)).All(&evns)
	if err != nil {
		return domain.Events{}, err
	}

	count, err := query.TotalEntries()
	if err != nil {
		return domain.Events{}, err
	}

	pages, err := paginate.TotalPages()
	if err != nil {
		return domain.Events{}, err
	}

	result := domain.Events{
		Items: r.mapModelToDomainCollection(evns),
		Pages: uint64(pages),
		Total: count,
	}

	return result, nil
}

func (r EventRepository) Update(t domain.Event) (domain.Event, error) {
	evn := r.mapDomainToModel(t)
	evn.UpdatedDate = time.Now()
	err := r.coll.Find(db.Cond{"id": evn.Id, "deleted_date": nil}).Update(&evn)
	if err != nil {
		return domain.Event{}, err
	}
	return r.mapModelToDomain(evn), nil
}

func (r EventRepository) Delete(id uint64) error {
	return r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

func (r EventRepository) mapDomainToModel(d domain.Event) event {
	return event{
		Id:          d.Id,
		User_Id:     d.User_Id,
		Title:       d.Title,
		Description: d.Description,
		Date:        d.Date,
		Image:       d.Image,
		City:        d.City,
		Location:    d.Location,
		Lat:         d.Lat,
		Lon:         d.Lon,
		CreatedDate: d.CreatedDate,
		UpdatedDate: d.UpdatedDate,
		DeletedDate: d.DeletedDate,
	}
}

func (r EventRepository) mapModelToDomain(m event) domain.Event {
	return domain.Event{
		Id:          m.Id,
		User_Id:     m.User_Id,
		Title:       m.Title,
		Description: m.Description,
		Date:        m.Date,
		Image:       m.Image,
		City:        m.City,
		Location:    m.Location,
		Lat:         m.Lat,
		Lon:         m.Lon,
		CreatedDate: m.CreatedDate,
		UpdatedDate: m.UpdatedDate,
		DeletedDate: m.DeletedDate,
	}
}

func (r EventRepository) mapModelToDomainCollection(evn []event) []domain.Event {
	var events []domain.Event
	for _, ev := range evn {
		events = append(events, r.mapModelToDomain(ev))
	}
	return events
}
