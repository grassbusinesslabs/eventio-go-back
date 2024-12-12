package app

import (
	"log"

	"github.com/grassbusinesslabs/eventio-go-back/internal/domain"
	"github.com/grassbusinesslabs/eventio-go-back/internal/infra/database"
)

type EventService struct {
	eventRepo database.EventRepository
}

func NewEventService(er database.EventRepository) EventService {
	return EventService{
		eventRepo: er,
	}
}

func (s EventService) Save(t domain.Event) (domain.Event, error) {
	evn, err := s.eventRepo.Save(t)
	if err != nil {
		log.Printf("EventService -> Save -> s.eventRepo.Save: %s", err)
		return domain.Event{}, err
	}
	return evn, nil
}

func (s EventService) Find(id uint64) (domain.Event, error) {
	evn, err := s.eventRepo.Find(id)
	if err != nil {
		log.Printf("EventService -> Find -> s.eventRepo.Find: %s", err)
		return domain.Event{}, err
	}
	return evn, nil
}

func (s EventService) FindList() ([]domain.Event, error) {
	events, err := s.eventRepo.FindList()
	if err != nil {
		log.Printf("EventService -> FindList -> s.eventRepo.FindList: %s", err)
		return nil, err
	}
	return events, nil
}

func (s EventService) Update(t domain.Event) (domain.Event, error) {
	e, err := s.eventRepo.Update(t)
	if err != nil {
		log.Printf("EventService -> Update -> s.eventRepo.Update: %s", err)
		return domain.Event{}, err
	}
	return e, nil
}

func (s EventService) Delete(id uint64) error {
	err := s.eventRepo.Delete(id)
	if err != nil {
		log.Printf("EventService: %s", err)
		return err
	}

	return nil
}
