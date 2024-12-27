package requests

import (
	"github.com/grassbusinesslabs/eventio-go-back/internal/domain"
)

type SubscriptionRequest struct {
	Event_Id int64 `json:"event_id" validate:"required"`
}

func (r SubscriptionRequest) ToDomainModel() (interface{}, error) {
	return domain.Subscription{
		Event_Id: uint64(r.Event_Id),
	}, nil
}
