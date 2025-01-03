package app

import (
	"log"

	"github.com/grassbusinesslabs/eventio-go-back/internal/domain"
	"github.com/grassbusinesslabs/eventio-go-back/internal/infra/database"
)

type SubscriptionService interface {
	Save(t domain.Subscription) (domain.Subscription, error)
	Delete(t domain.Subscription) error
	CountByEvent(event_Id uint64) (uint64, error)
	GetUserSubsId(user_Id uint64, page domain.Pagination) ([]uint64, domain.Page, error)
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

func (s subscriptionService) CountByEvent(event_Id uint64) (uint64, error) {
	count, err := s.subsRepo.CountByEvent(event_Id)
	if err != nil {
		log.Printf("subscriptionService: %s", err)
		return count, err
	}

	return count, err
}

func (s subscriptionService) GetUserSubsId(user_Id uint64, page domain.Pagination) ([]uint64, domain.Page, error) {
	subsId, total, err := s.subsRepo.GetUserSubsId(user_Id, page)
	if err != nil {
		log.Printf("subscriptionService: %s", err)
		return subsId, total, err
	}

	return subsId, total, err
}
