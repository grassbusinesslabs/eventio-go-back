package app

import (
	"log"
	"time"

	"github.com/grassbusinesslabs/eventio-go-back/internal/domain"
	"github.com/grassbusinesslabs/eventio-go-back/internal/infra/database"
)

type EventService interface {
	Save(event domain.Event) (domain.Event, error)
	Find(id uint64) (domain.Event, error)
	FindList(city string) ([]domain.Event, error)
	FindListByUser(id uint64) ([]domain.Event, error)
	FindListByDay(time time.Time, city string) ([]domain.Event, error)
	FindListByMonth(time time.Time, city string) ([]domain.Event, error)
	FindListBySearch(title string, city string) ([]domain.Event, error)
	Update(event domain.Event) (domain.Event, error)
	Delete(id uint64) error
}

type eventService struct {
	eventRepo database.EventRepository
}

func NewEventService(er database.EventRepository) EventService {
	return eventService{
		eventRepo: er,
	}
}

func (s eventService) Save(t domain.Event) (domain.Event, error) {
	evn, err := s.eventRepo.Save(t)
	if err != nil {
		log.Printf("EventService -> Save -> s.eventRepo.Save: %s", err)
		return domain.Event{}, err
	}
	return evn, nil
}

func (s eventService) Find(id uint64) (domain.Event, error) {
	evn, err := s.eventRepo.Find(id)
	if err != nil {
		log.Printf("EventService -> Find -> s.eventRepo.Find: %s", err)
		return domain.Event{}, err
	}
	return evn, nil
}

func (s eventService) FindList(city string) ([]domain.Event, error) {
	events, err := s.eventRepo.FindList(city)
	if err != nil {
		log.Printf("EventService -> FindList -> s.eventRepo.FindList: %s", err)
		return nil, err
	}
	return events, nil
}

func (s eventService) FindListByUser(id uint64) ([]domain.Event, error) {
	events, err := s.eventRepo.FindListByUser(id)
	if err != nil {
		log.Printf("EventService -> FindListByUser -> s.eventRepo.FindListByUser: %s", err)
		return nil, err
	}
	return events, nil
}

func (s eventService) FindListByDay(date time.Time, city string) ([]domain.Event, error) {
	events, err := s.eventRepo.FindListByDay(date, city)
	if err != nil {
		log.Printf("EventService -> FindListByDay -> s.eventRepo.FindListByDay: %s", err)
		return nil, err
	}
	return events, nil
}

func (s eventService) FindListByMonth(date time.Time, city string) ([]domain.Event, error) {
	events, err := s.eventRepo.FindListByMonth(date, city)
	if err != nil {
		log.Printf("EventService -> FindListByMonth -> s.eventRepo.FindListByMonth: %s", err)
		return nil, err
	}
	return events, nil
}

func (s eventService) FindListBySearch(search string, city string) ([]domain.Event, error) {
	events, err := s.eventRepo.FindListBySearch(search, city)
	if err != nil {
		log.Printf("EventService -> FindListBySearch -> s.eventRepo.FindListBySearch: %s", err)
		return nil, err
	}
	return events, nil
}

func (s eventService) Update(t domain.Event) (domain.Event, error) {
	e, err := s.eventRepo.Update(t)
	if err != nil {
		log.Printf("EventService -> Update -> s.eventRepo.Update: %s", err)
		return domain.Event{}, err
	}
	return e, nil
}

func (s eventService) Delete(id uint64) error {
	err := s.eventRepo.Delete(id)
	if err != nil {
		log.Printf("EventService: %s", err)
		return err
	}

	return nil
}
