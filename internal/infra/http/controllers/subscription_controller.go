package controllers

import (
	"log"
	"net/http"
	"strconv"

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

		pageStr := r.URL.Query().Get("page")
		page, err := strconv.ParseUint(pageStr, 10, 64)
		if err != nil {
			log.Printf("EventController -> strconv.ParseUint: %s", err)
			BadRequest(w, err)
			return
		}

		pagination := domain.Pagination{
			Page:         page,
			CountPerPage: 15,
		}

		subsId, total, err := c.subscriptionService.GetUserSubsId(user.Id, pagination)
		if err != nil {
			log.Printf("SubscriptionController -> FindUserSubs -> GetUserSubsId: %s", err)
			InternalServerError(w, err)
			return
		}

		var result domain.Events
		for _, e := range subsId {
			event, err := c.eventService.FindById(e)
			if err != nil {
				log.Printf("SubscriptionController -> FindUserSubs -> FindById: %s", err)
				InternalServerError(w, err)
				return
			}
			result.Items = append(result.Items, event)
		}

		var eventsDto resources.EventsDto
		eventsDto = eventsDto.DomainToDto(result)
		eventsDto.Pages = uint64((total + pagination.CountPerPage - 1) / pagination.CountPerPage)
		eventsDto.Total = total

		Success(w, eventsDto)
	}
}
