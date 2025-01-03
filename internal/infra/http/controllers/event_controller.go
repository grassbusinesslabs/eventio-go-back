package controllers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/grassbusinesslabs/eventio-go-back/internal/app"
	"github.com/grassbusinesslabs/eventio-go-back/internal/domain"
	"github.com/grassbusinesslabs/eventio-go-back/internal/infra/database"
	"github.com/grassbusinesslabs/eventio-go-back/internal/infra/filesystem"
	"github.com/grassbusinesslabs/eventio-go-back/internal/infra/http/requests"
	"github.com/grassbusinesslabs/eventio-go-back/internal/infra/http/resources"
)

type EventController struct {
	eventService app.EventService
	subsService  app.SubscriptionService
	imageStorage filesystem.ImageStorageService
}

func NewEventController(ev app.EventService, sub app.SubscriptionService, imgStorage filesystem.ImageStorageService) EventController {
	return EventController{
		eventService: ev,
		subsService:  sub,
		imageStorage: imgStorage,
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

func (c EventController) FindById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		IdStr := r.URL.Query().Get("id")

		Id, err := strconv.ParseUint(IdStr, 10, 64)
		if err != nil {
			log.Printf("EventController -> strconv.ParseUint: %s", err)
			BadRequest(w, err)
			return
		}

		event, err := c.eventService.FindById(Id)
		if err != nil {
			log.Printf("EventController -> FindList -> c.eventService.FindList: %s", err)
			InternalServerError(w, err)
			return
		}

		var eventDto resources.EventDto
		eventDto = eventDto.DomainToDto(event)
		Success(w, eventDto)
	}
}

func (c EventController) FindListBy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)

		city := r.URL.Query().Get("city")
		search := r.URL.Query().Get("search")
		dayunix := r.URL.Query().Get("day")
		monthunix := r.URL.Query().Get("month")
		yearunix := r.URL.Query().Get("year")
		location := r.URL.Query().Get("location")
		userStr := r.URL.Query().Get("user")
		pageStr := r.URL.Query().Get("page")
		page, err := strconv.ParseUint(pageStr, 10, 64)
		if err != nil {
			log.Printf("EventController -> strconv.ParseUint: %s", err)
			BadRequest(w, err)
			return
		}

		var str database.EventSearchParams

		if userStr != "" {
			str.User_Id = user.Id
		}

		str.City = city
		str.Search = search
		str.Location = location
		str.Pagination = domain.Pagination{
			Page:         page,
			CountPerPage: 15,
		}

		if dayunix != "" {
			dayunix, err := strconv.ParseUint(dayunix, 10, 64)
			if err != nil {
				log.Printf("EventController -> strconv.ParseUint: %s", err)
				BadRequest(w, err)
				return
			}
			date := time.Unix(int64(dayunix), 0)

			str.DateDay = &date

			events, err := c.eventService.FindListBy(str)
			if err != nil {
				log.Printf("EventController -> FindList -> c.eventService.FindList: %s", err)
				InternalServerError(w, err)
				return
			}

			var eventsDto resources.EventsDto
			eventsDto = eventsDto.DomainToDto(events)
			Success(w, eventsDto)
		} else if monthunix != "" {
			monthunix, err := strconv.ParseUint(monthunix, 10, 64)
			if err != nil {
				log.Printf("EventController -> strconv.ParseUint: %s", err)
				BadRequest(w, err)
				return
			}
			date := time.Unix(int64(monthunix), 0)

			str.DateMonth = &date

			events, err := c.eventService.FindListBy(str)
			if err != nil {
				log.Printf("EventController -> FindList -> c.eventService.FindList: %s", err)
				InternalServerError(w, err)
				return
			}

			var eventsDto resources.EventsDto
			eventsDto = eventsDto.DomainToDto(events)
			Success(w, eventsDto)
		} else if yearunix != "" {
			yearunix, err := strconv.ParseUint(yearunix, 10, 64)
			if err != nil {
				log.Printf("EventController -> strconv.ParseUint: %s", err)
				BadRequest(w, err)
				return
			}
			date := time.Unix(int64(yearunix), 0)

			str.DateYear = &date

			events, err := c.eventService.FindListBy(str)
			if err != nil {
				log.Printf("EventController -> FindList -> c.eventService.FindList: %s", err)
				InternalServerError(w, err)
				return
			}

			var eventsDto resources.EventsDto
			eventsDto = eventsDto.DomainToDto(events)
			Success(w, eventsDto)
		} else {
			events, err := c.eventService.FindListBy(str)
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
		event.City = updateEvent.City
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

func (c EventController) UploadImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, _, err := r.FormFile("image")
		if err != nil {
			log.Printf("EventController -> UploadImage -> FormFile: %s", err)
			BadRequest(w, err)
			return
		}
		defer file.Close()

		fileContent, err := io.ReadAll(file)
		if err != nil {
			log.Printf("EventController -> UploadImage -> ReadAll: %s", err)
			InternalServerError(w, err)
			return
		}

		event := r.Context().Value(EventKey).(domain.Event)

		filename := fmt.Sprintf("%d.png", event.Id)

		filePath := fmt.Sprintf("file_storage/%s", filename)
		_, err = os.Stat(filePath)
		if os.IsNotExist(err) {
			err = c.imageStorage.SaveImage(filename, fileContent)
			if err != nil {
				log.Printf("EventController -> UploadImage -> SaveImage: %s", err)
				InternalServerError(w, err)
				return
			}

			event.Image = filename
			event, err = c.eventService.Update(event)
			if err != nil {
				log.Printf("EventController -> UploadImage -> Update: %s", err)
				InternalServerError(w, err)
				return
			}

			Success(w, map[string]string{"message": "File saved!", "path": filename})
		} else {
			err := c.imageStorage.DeleteImage(filename)
			if err != nil {
				log.Printf("EventController -> DeleteImage -> DeleteImage: %s", err)
				InternalServerError(w, err)
				return
			}

			err = c.imageStorage.SaveImage(filename, fileContent)
			if err != nil {
				log.Printf("EventController -> UploadImage -> SaveImage: %s", err)
				InternalServerError(w, err)
				return
			}

			event.Image = filename
			event, err = c.eventService.Update(event)
			if err != nil {
				log.Printf("EventController -> UploadImage -> Update: %s", err)
				InternalServerError(w, err)
				return
			}

			Success(w, map[string]string{"message": "File saved!", "path": filename})
		}
	}
}

func (c EventController) DeleteImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		event := r.Context().Value(EventKey).(domain.Event)

		filename := fmt.Sprintf("%d.png", event.Id)
		err := c.imageStorage.DeleteImage(filename)
		if err != nil {
			log.Printf("EventController -> DeleteImage -> DeleteImage: %s", err)
			InternalServerError(w, err)
			return
		}
		event.Image = ""
		event, err = c.eventService.Update(event)
		if err != nil {
			log.Printf("EventController -> DeleteImage -> Update: %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, map[string]string{"message": "File deleted!"})
	}
}
