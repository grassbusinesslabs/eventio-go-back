package controllers

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/grassbusinesslabs/eventio-go-back/internal/app"
	"github.com/grassbusinesslabs/eventio-go-back/internal/domain"
	"github.com/grassbusinesslabs/eventio-go-back/internal/infra/filesystem"
	"github.com/grassbusinesslabs/eventio-go-back/internal/infra/http/requests"
	"github.com/grassbusinesslabs/eventio-go-back/internal/infra/http/resources"
)

type UserController struct {
	userService  app.UserService
	authService  app.AuthService
	imageStorage filesystem.ImageStorageService
}

func NewUserController(us app.UserService, as app.AuthService, imgStorage filesystem.ImageStorageService) UserController {
	return UserController{
		userService:  us,
		authService:  as,
		imageStorage: imgStorage,
	}
}

func (c UserController) FindMe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		Success(w, resources.UserDto{}.DomainToDto(user))
	}
}

func (c UserController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := requests.Bind(r, requests.UpdateUserRequest{}, domain.User{})
		if err != nil {
			log.Printf("UserController: %s", err)
			BadRequest(w, err)
			return
		}

		u := r.Context().Value(UserKey).(domain.User)
		u.FirstName = user.FirstName
		u.SecondName = user.SecondName
		u.Email = user.Email
		user, err = c.userService.Update(u)
		if err != nil {
			log.Printf("UserController: %s", err)
			InternalServerError(w, err)
			return
		}

		var userDto resources.UserDto
		Success(w, userDto.DomainToDto(user))
	}
}

func (c UserController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(UserKey).(domain.User)

		err := c.userService.Delete(u.Id)
		if err != nil {
			log.Printf("UserController: %s", err)
			InternalServerError(w, err)
			return
		}

		Ok(w)
	}
}

func (c UserController) UploadUserImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, _, err := r.FormFile("image")
		if err != nil {
			log.Printf("EventController -> UploadUserImage -> FormFile: %s", err)
			BadRequest(w, err)
			return
		}
		defer file.Close()

		fileContent, err := io.ReadAll(file)
		if err != nil {
			log.Printf("EventController -> UploadUserImage -> ReadAll: %s", err)
			InternalServerError(w, err)
			return
		}

		user := r.Context().Value(UserKey).(domain.User)

		filename := fmt.Sprintf("/user_image/%d.png", user.Id)
		err = c.imageStorage.SaveImage(filename, fileContent)
		if err != nil {
			log.Printf("EventController -> UploadUserImage -> SaveImage: %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, map[string]string{"message": "File saved!", "path": filename})
	}
}

func (c UserController) DeleteUserImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)

		filename := fmt.Sprintf("/user_image/%d.png", user.Id)
		err := c.imageStorage.DeleteImage(filename)
		if err != nil {
			log.Printf("EventController -> DeleteUserImage -> DeleteImage: %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, map[string]string{"message": "File deleted!"})
	}
}

func (c UserController) UpdateUserImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)

		filename := fmt.Sprintf("/user_image/%d.png", user.Id)
		err := c.imageStorage.DeleteImage(filename)
		if err != nil {
			log.Printf("EventController -> UpdateUserImage -> DeleteImage: %s", err)
			InternalServerError(w, err)
			return
		}

		file, _, err := r.FormFile("image")
		if err != nil {
			log.Printf("EventController -> UpdateUserImage -> FormFile: %s", err)
			BadRequest(w, err)
			return
		}
		defer file.Close()

		fileContent, err := io.ReadAll(file)
		if err != nil {
			log.Printf("EventController -> UpdateUserImage -> ReadAll: %s", err)
			InternalServerError(w, err)
			return
		}

		err = c.imageStorage.SaveImage(filename, fileContent)
		if err != nil {
			log.Printf("EventController -> UpdateUserImage -> SaveImage: %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, map[string]string{"message": "File updated!", "path": filename})
	}
}
