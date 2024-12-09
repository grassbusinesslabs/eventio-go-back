package database

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/upper/db/v4"
)

const EventsTableName = "events"

type event struct {
	EventId     uint64     `db:"eventid,omitempty"`
	UserId      uint64     `db:"userid,omitempty"`
	Tytle       string     `db:"tytle"`
	Description string     `db:"description"`
	Date        time.Time  `db:"date"`
	Image       string     `db:"image"`
	Location    string     `db:"location"`
	Lat         float64    `db:"lat"`
	Lon         float64    `db:"lon"`
	CreatedDate time.Time  `db:"created_date,omitempty"`
	UpdatedDate time.Time  `db:"updated_date,omitempty"`
	DeletedDate *time.Time `db:"deleted_date,omitempty"`
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

func (r EventRepository) Find(eventid uint64) (domain.Event, error) {
	var evn event

	err := r.coll.Find(db.Cond{"eventid": eventid}).One(&evn)
	if err != nil {
		return domain.Event{}, err
	}

	ev := r.mapModelToDomain(evn)
	return ev, nil
}

func (r EventRepository) FindList() ([]domain.Event, error) {
	var evns []event

	err := r.coll.Find().All(&evns)
	if err != nil {
		return nil, err
	}

	evs := r.mapModelToDomainCollection(evns)
	return evs, nil
}

func (r EventRepository) Update(t domain.Event) (domain.Event, error) {
	evn := r.mapDomainToModel(t)
	evn.UpdatedDate = time.Now()
	err := r.coll.Find(db.Cond{"eventid": evn.EventId, "deleted_date": nil}).Update(&evn)
	if err != nil {
		return domain.Event{}, err
	}
	return r.mapModelToDomain(evn), nil
}

func (r EventRepository) Delete(eventid uint64) error {
	return r.coll.Find(db.Cond{"eventid": eventid, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

func (r EventRepository) mapDomainToModel(d domain.Event) event {
	return event{
		EventId:     d.EventId,
		UserId:      d.UserId,
		Tytle:       d.Tytle,
		Description: d.Description,
		Date:        d.Date,
		Image:       d.Image,
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
		EventId:     m.EventId,
		UserId:      m.UserId,
		Tytle:       m.Tytle,
		Description: m.Description,
		Date:        m.Date,
		Image:       m.Image,
		Location:    m.Location,
		Lat:         m.Lat,
		Lon:         m.Lon,
		CreatedDate: m.CreatedDate,
		UpdatedDate: m.UpdatedDate,
		DeletedDate: m.DeletedDate,
	}
}

func (r EventRepository) mapModelToDomainCollection(evn []event) []domain.Event {
	var evns []domain.Event
	for _, ev := range evn {
		evns = append(evns, r.mapModelToDomain(ev))
	}
	return evns
}
