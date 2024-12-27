package app

import (
	"log"

	"github.com/grassbusinesslabs/eventio-go-back/internal/domain"
	"github.com/grassbusinesslabs/eventio-go-back/internal/infra/database"
)

type SubscriptionService interface {
	Save(t domain.Subscription) (domain.Subscription, error)
	Delete(t domain.Subscription) error
}

type subscriptionService struct {
	subsRepo database.SubscriptionRepository
}

func NewSubscriptionService(er database.SubscriptionRepository) SubscriptionService {
	return subscriptionService{
		subsRepo: er,
	}
}

func (s subscriptionService) Save(t domain.Subscription) (domain.Subscription, error) {
	evn, err := s.subsRepo.Save(t)
	if err != nil {
		log.Printf("subscriptionService -> Save -> s.subsRepo.Save: %s", err)
		return domain.Subscription{}, err
	}
	return evn, nil
}

func (s subscriptionService) Delete(t domain.Subscription) error {
	err := s.subsRepo.Delete(t)
	if err != nil {
		log.Printf("subscriptionService: %s", err)
		return err
	}

	return nil
}
