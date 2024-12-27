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
