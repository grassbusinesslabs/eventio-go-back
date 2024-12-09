package app

import (
	"log"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/database"
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
	evns, err := s.eventRepo.FindList()
	if err != nil {
		log.Printf("EventService -> FindList -> s.eventRepo.FindList: %s", err)
		return nil, err
	}
	return evns, nil
}

func (s EventService) Update(t domain.Event) (domain.Event, error) {
	evns, err := s.eventRepo.Update(t)
	if err != nil {
		log.Printf("EventService -> Update -> s.eventRepo.Update: %s", err)
		return domain.Event{}, err
	}
	return evns, nil
}

func (s EventService) Delete(eventid uint64) error {
	err := s.eventRepo.Delete(eventid)
	if err != nil {
		log.Printf("EventService: %s", err)
		return err
	}

	return nil
}
