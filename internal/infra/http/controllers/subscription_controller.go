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
}

func NewSubscriptionController(ev app.SubscriptionService) SubscriptionController {
	return SubscriptionController{
		subscriptionService: ev,
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
