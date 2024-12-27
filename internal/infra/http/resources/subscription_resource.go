package resources

import (
	"github.com/grassbusinesslabs/eventio-go-back/internal/domain"
)

type SubscriptionDto struct {
	User_Id  uint64 `json:"user_id"`
	Event_Id uint64 `json:"event_id"`
}

func (d SubscriptionDto) DomainToDto(t domain.Subscription) SubscriptionDto {
	return SubscriptionDto{
		User_Id:  t.User_Id,
		Event_Id: t.Event_Id,
	}
}
