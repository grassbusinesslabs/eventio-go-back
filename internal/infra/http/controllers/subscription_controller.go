package controllers

import (
	"log"
	"net/http"

	"github.com/grassbusinesslabs/eventio-go-back/internal/app"
	"github.com/grassbusinesslabs/eventio-go-back/internal/domain"
	"github.com/grassbusinesslabs/eventio-go-back/internal/infra/http/requests"
	"github.com/grassbusinesslabs/eventio-go-back/internal/infra/http/resources"
)

type SubscriptionController struct {
	subscriptionService app.SubscriptionService
	eventService        app.EventService
}

func NewSubscriptionController(subs app.SubscriptionService, ev app.EventService) SubscriptionController {
	return SubscriptionController{
		subscriptionService: subs,
		eventService:        ev,
	}
}

func (c SubscriptionController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subs, err := requests.Bind(r, requests.SubscriptionRequest{}, domain.Subscription{})
		if err != nil {
			log.Printf("SubscriptionController -> Save -> requests.Bind: %s", err)
			BadRequest(w, err)
			return
		}

		user := r.Context().Value(UserKey).(domain.User)
		subs.User_Id = user.Id

		subs, err = c.subscriptionService.Save(subs)
		if err != nil {
			log.Printf("EventController -> Save -> c.eventService.Save: %s", err)
			InternalServerError(w, err)
			return
		}

		var subsDto resources.SubscriptionDto
		subsDto = subsDto.DomainToDto(subs)
		Created(w, subsDto)
	}
}

func (c SubscriptionController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subs, err := requests.Bind(r, requests.SubscriptionRequest{}, domain.Subscription{})
		if err != nil {
			log.Printf("SubscriptionController -> Save -> requests.Bind: %s", err)
			BadRequest(w, err)
			return
		}

		user := r.Context().Value(UserKey).(domain.User)
		subs.User_Id = user.Id

		err = c.subscriptionService.Delete(subs)
		if err != nil {
			log.Printf("EventController: %s", err)
			InternalServerError(w, err)
			return
		}

		Ok(w)
	}
}

func (c SubscriptionController) FindUserSubs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)

		subsId, err := c.subscriptionService.GetUserSubsId(user.Id)
		if err != nil {
			log.Printf("EventController -> Find -> c.subsService.CountByEvent: %s", err)
			InternalServerError(w, err)
			return
		}

		events := make([]domain.Event, len(subsId))
		for i, e := range subsId {
			events[i], err = c.eventService.Find(e)
			if err != nil {
				log.Printf("EventController: %s", err)
				return
			}
		}

		var eventsDto resources.EventsDto
		eventsDto = eventsDto.DomainToDto(events)
		Success(w, eventsDto)
	}
}
