package controllers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/grassbusinesslabs/eventio-go-back/internal/app"
	"github.com/grassbusinesslabs/eventio-go-back/internal/domain"
	"github.com/grassbusinesslabs/eventio-go-back/internal/infra/http/requests"
	"github.com/grassbusinesslabs/eventio-go-back/internal/infra/http/resources"
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
		event.User_Id = user.Id

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
		event := r.Context().Value(EventKey).(domain.Event)

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

func (c EventController) FindListByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		events, err := c.eventService.FindListByUser(user.Id)
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

func (c EventController) FindListByDate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dateunix, err := strconv.ParseUint(r.URL.Query().Get("date"), 10, 64)
		if err != nil {
			log.Printf("EventController -> strconv.ParseUint: %s", err)
			BadRequest(w, err)
			return
		}
		date := time.Unix(int64(dateunix), 0)

		events, err := c.eventService.FindListByDate(date)
		if err != nil {
			log.Printf("EventController -> FindListByDate -> FindListByDate: %s", err)
			InternalServerError(w, err)
			return
		}
		var eventsDto resources.EventsDto
		eventsDto = eventsDto.DomainToDto(events)
		Success(w, eventsDto)
	}
}

func (c EventController) FindListByTitle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Query().Get("title")

		events, err := c.eventService.FindListByTitle(title)
		if err != nil {
			log.Printf("EventController -> FindListByTitle -> FindListByTitle: %s", err)
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
			log.Printf("EventController -> Update -> requests.Bind: %s", err)
			BadRequest(w, err)
			return
		}
		event := r.Context().Value(EventKey).(domain.Event)

		event.Title = updateEvent.Title
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
		event := r.Context().Value(EventKey).(domain.Event)

		err := c.eventService.Delete(event.Id)
		if err != nil {
			log.Printf("EventController: %s", err)
			InternalServerError(w, err)
			return
		}

		Ok(w)
	}
}
