package controllers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/app"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/requests"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/resources"
	"github.com/go-chi/chi/v5"
)

type EventController struct {
	eventService app.EventService
}

func NewEventController(ev app.EventService) EventController {
	return EventController{
		eventService: ev,
	}
}

func (c EventController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		event, err := requests.Bind(r, requests.EventRequest{}, domain.Event{})
		if err != nil {
			log.Printf("EventController -> Save -> requests.Bind: %s", err)
			BadRequest(w, err)
			return
		}

		user := r.Context().Value(UserKey).(domain.User)
		event.UserId = user.Id

		event, err = c.eventService.Save(event)
		if err != nil {
			log.Printf("EventController -> Save -> c.eventService.Save: %s", err)
			InternalServerError(w, err)
			return
		}

		var eventDto resources.EventDto
		eventDto = eventDto.DomainToDto(event)
		Created(w, eventDto)
	}
}

func (c EventController) Find() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventId, err := strconv.ParseUint(chi.URLParam(r, "eventid"), 10, 64)
		if err != nil {
			log.Printf("EventController -> Find -> strconv.ParseUint: %s", err)
			BadRequest(w, err)
			return
		}

		event, err := c.eventService.Find(eventId)
		if err != nil {
			log.Printf("EventController -> Find -> c.eventService.Find: %s", err)
			InternalServerError(w, err)
			return
		}

		var eventDto resources.EventDto
		eventDto = eventDto.DomainToDto(event)
		Success(w, eventDto)
	}
}

func (c EventController) FindList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		events, err := c.eventService.FindList()
		if err != nil {
			log.Printf("EventController -> FindList -> c.eventService.FindList: %s", err)
			InternalServerError(w, err)
			return
		}

		var eventsDto resources.EventsDto
		eventsDto = eventsDto.DomainToDto(events)
		Success(w, eventsDto)
	}
}

func (c EventController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		updateEvent, err := requests.Bind(r, requests.EventRequest{}, domain.Event{})
		if err != nil {
			log.Printf("EventController: %s", err)
			BadRequest(w, err)
			return
		}

		eventId, err := strconv.ParseUint(chi.URLParam(r, "eventid"), 10, 64)
		if err != nil {
			log.Printf("EventController -> Update -> strconv.ParseUint: %s", err)
			BadRequest(w, err)
			return
		}

		user := r.Context().Value(UserKey).(domain.User)
		event, err := c.eventService.Find(eventId)
		if err != nil {
			log.Printf("EventController -> Update -> c.eventService.Find: %s", err)
			InternalServerError(w, err)
			return
		}
		if user.Id != event.UserId {
			err = errors.New("Access denied!")
			Forbidden(w, err)
			return
		}

		event.Tytle = updateEvent.Tytle
		event.Description = updateEvent.Description
		event.Date = updateEvent.Date
		event.Image = updateEvent.Image
		event.Location = updateEvent.Location
		event.Lat = updateEvent.Lat
		event.Lon = updateEvent.Lon
		event, err = c.eventService.Update(event)
		if err != nil {
			log.Printf("EventController: %s", err)
			InternalServerError(w, err)
			return
		}

		var eventDto resources.EventDto
		Success(w, eventDto.DomainToDto(event))
	}
}

func (c EventController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventId, err := strconv.ParseUint(chi.URLParam(r, "eventid"), 10, 64)
		if err != nil {
			log.Printf("EventController -> Delete -> strconv.ParseUint: %s", err)
			BadRequest(w, err)
			return
		}

		user := r.Context().Value(UserKey).(domain.User)
		event, err := c.eventService.Find(eventId)
		if err != nil {
			log.Printf("EventController -> Delete -> c.eventService.Find: %s", err)
			InternalServerError(w, err)
			return
		}
		if user.Id != event.UserId {
			err = errors.New("Access denied!")
			Forbidden(w, err)
			return
		}
		//e := r.Context().Value(eventId).(domain.Event)

		err = c.eventService.Delete(eventId)
		if err != nil {
			log.Printf("EventController: %s", err)
			InternalServerError(w, err)
			return
		}

		Ok(w)
	}
}
