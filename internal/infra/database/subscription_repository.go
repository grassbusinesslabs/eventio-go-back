package database

import (
	"log"

	"github.com/grassbusinesslabs/eventio-go-back/internal/domain"
	"github.com/upper/db/v4"
)

const SubscriptionsTableName = "subscriptions"

type subscription struct {
	Event_Id uint64 `db:"event_id"`
	User_Id  uint64 `db:"user_id"`
}

type SubscriptionRepository struct {
	coll db.Collection
	sess db.Session
}

func NewSubscrRepository(sess db.Session) SubscriptionRepository {
	return SubscriptionRepository{
		coll: sess.Collection(SubscriptionsTableName),
		sess: sess,
	}
}

func (r SubscriptionRepository) Save(t domain.Subscription) (domain.Subscription, error) {
	evn := r.mapDomainToModel(t)

	err := r.coll.InsertReturning(&evn)
	if err != nil {
		return domain.Subscription{}, err
	}

	t = r.mapModelToDomain(evn)
	return t, nil
}

func (r SubscriptionRepository) Delete(t domain.Subscription) error {
	delete := r.coll.Find(db.Cond{"event_id": t.Event_Id, "user_id": t.User_Id}).Delete()
	if delete != nil {
		log.Printf("SubscriptionRepository: %s", delete)
		return delete
	}
	return delete
}

func (r SubscriptionRepository) CountByEvent(event_Id uint64) (uint64, error) {
	count, err := r.coll.Find(db.Cond{"event_id": event_Id}).Count()
	if err != nil {
		log.Printf("SubscriptionRepository -> CountByEvent: %v", err)
		return 0, err
	}
	return count, nil
}

func (r SubscriptionRepository) GetUserSubsId(user_Id uint64, p domain.Pagination) ([]uint64, domain.Page, error) {
	var subscriptions []subscription

	res := r.coll.Find(db.Cond{"user_id": user_Id})

	paginate := res.Paginate(uint(p.CountPerPage))
	err := paginate.Page(uint(p.Page)).All(&subscriptions)
	if err != nil {
		log.Printf("SubscriptionRepository -> GetUserSubsId -> Paginate: %v", err)
		return nil, domain.Page{}, err
	}

	total, err := res.TotalEntries()
	if err != nil {
		log.Printf("SubscriptionRepository -> GetUserSubsId -> TotalEntries: %v", err)
		return nil, domain.Page{}, err
	}

	pages, err := paginate.TotalPages()
	if err != nil {
		log.Printf("SubscriptionRepository -> GetUserSubsId -> TotalPages: %v", err)
		return nil, domain.Page{}, err
	}

	var eventsId []uint64
	for _, subsc := range subscriptions {
		eventsId = append(eventsId, subsc.Event_Id)
	}

	page := domain.Page{
		Pages: uint64(pages),
		Total: total,
	}

	return eventsId, page, nil
}

func (r SubscriptionRepository) mapDomainToModel(d domain.Subscription) subscription {
	return subscription{
		Event_Id: d.Event_Id,
		User_Id:  d.User_Id,
	}
}

func (r SubscriptionRepository) mapModelToDomain(m subscription) domain.Subscription {
	return domain.Subscription{
		Event_Id: m.Event_Id,
		User_Id:  m.User_Id,
	}
}
