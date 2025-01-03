package container

import (
	"log"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/grassbusinesslabs/eventio-go-back/config"
	"github.com/grassbusinesslabs/eventio-go-back/internal/app"
	"github.com/grassbusinesslabs/eventio-go-back/internal/infra/database"
	"github.com/grassbusinesslabs/eventio-go-back/internal/infra/filesystem"
	"github.com/grassbusinesslabs/eventio-go-back/internal/infra/http/controllers"
	"github.com/grassbusinesslabs/eventio-go-back/internal/infra/http/middlewares"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
)

type Container struct {
	Middlewares
	Services
	Controllers
}

type Middlewares struct {
	AuthMw  func(http.Handler) http.Handler
	EventMw func(http.Handler) http.Handler
}

type Services struct {
	app.AuthService
	app.UserService
	app.EventService
	app.SubscriptionService
}

type Controllers struct {
	AuthController         controllers.AuthController
	UserController         controllers.UserController
	EventController        controllers.EventController
	SubscriptionController controllers.SubscriptionController
}

func New(conf config.Configuration) Container {
	tknAuth := jwtauth.New("HS256", []byte(conf.JwtSecret), nil)
	sess := getDbSess(conf)

	sessionRepository := database.NewSessRepository(sess)
	userRepository := database.NewUserRepository(sess)
	eventRepository := database.NewEventRepository(sess)
	subscriptionRepository := database.NewSubscrRepository(sess)

	userService := app.NewUserService(userRepository)
	authService := app.NewAuthService(sessionRepository, userRepository, tknAuth, conf.JwtTTL)
	eventService := app.NewEventService(eventRepository)
	subscriptionService := app.NewSubscriptionService(subscriptionRepository)

	imageStorage := filesystem.NewImageStorageService(conf)

	authController := controllers.NewAuthController(authService, userService)
	userController := controllers.NewUserController(userService, authService, imageStorage)
	eventController := controllers.NewEventController(eventService, subscriptionService, imageStorage)
	subscriptionController := controllers.NewSubscriptionController(subscriptionService, eventService)

	authMiddleware := middlewares.AuthMiddleware(tknAuth, authService, userService)
	eventMiddleware := middlewares.EventMiddleware(eventService)

	return Container{
		Middlewares: Middlewares{
			AuthMw:  authMiddleware,
			EventMw: eventMiddleware,
		},
		Services: Services{
			authService,
			userService,
			eventService,
			subscriptionService,
		},
		Controllers: Controllers{
			authController,
			userController,
			eventController,
			subscriptionController,
		},
	}
}

func getDbSess(conf config.Configuration) db.Session {
	sess, err := postgresql.Open(
		postgresql.ConnectionURL{
			User:     conf.DatabaseUser,
			Host:     conf.DatabaseHost,
			Password: conf.DatabasePassword,
			Database: conf.DatabaseName,
		})
	if err != nil {
		log.Fatalf("Unable to create new DB session: %q\n", err)
	}
	return sess
}
