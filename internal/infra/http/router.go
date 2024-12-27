package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/grassbusinesslabs/eventio-go-back/config"
	"github.com/grassbusinesslabs/eventio-go-back/config/container"
	"github.com/grassbusinesslabs/eventio-go-back/internal/domain"
	"github.com/grassbusinesslabs/eventio-go-back/internal/infra/http/controllers"
	"github.com/grassbusinesslabs/eventio-go-back/internal/infra/http/middlewares"

	"github.com/go-chi/chi/v5/middleware"
)

func Router(cont container.Container) http.Handler {

	router := chi.NewRouter()

	router.Use(middleware.RedirectSlashes, middleware.Logger, cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*", "http://localhost:5173/events", "capacitor://localhost"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Route("/api", func(apiRouter chi.Router) {
		// Health
		apiRouter.Route("/ping", func(healthRouter chi.Router) {
			healthRouter.Get("/", PingHandler())
			healthRouter.Handle("/*", NotFoundJSON())
		})

		apiRouter.Route("/v1", func(apiRouter chi.Router) {
			// Public routes
			apiRouter.Group(func(apiRouter chi.Router) {
				apiRouter.Route("/auth", func(apiRouter chi.Router) {
					AuthRouter(apiRouter, cont.AuthController, cont.AuthMw)
				})
			})

			// Protected routes
			apiRouter.Group(func(apiRouter chi.Router) {
				apiRouter.Use(cont.AuthMw)

				UserRouter(apiRouter, cont.UserController)
				EventRouter(apiRouter, cont.EventController, cont.EventMw)
				SubscriptionRouter(apiRouter, cont.SubscriptionController)
				apiRouter.Handle("/*", NotFoundJSON())
			})
		})
	})

	router.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Дозволити запити з будь-якого домену (або обмежити конкретним доменом)
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		workDir, _ := os.Getwd()
		filesDir := http.Dir(filepath.Join(workDir, config.GetConfiguration().FileStorageLocation))
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(filesDir))
		fs.ServeHTTP(w, r)
	})

	router.Options("/static/*", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173/events")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
	})

	return router
}

func AuthRouter(r chi.Router, ac controllers.AuthController, amw func(http.Handler) http.Handler) {
	r.Route("/", func(apiRouter chi.Router) {
		apiRouter.Post(
			"/register",
			ac.Register(),
		)
		apiRouter.Post(
			"/login",
			ac.Login(),
		)
		apiRouter.With(amw).Post(
			"/logout",
			ac.Logout(),
		)
	})
}

func UserRouter(r chi.Router, uc controllers.UserController) {
	uimwT := middlewares.UserImageMiddleware(true)
	uimwF := middlewares.UserImageMiddleware(false)
	r.Route("/users", func(apiRouter chi.Router) {
		apiRouter.Get(
			"/",
			uc.FindMe(),
		)
		apiRouter.Put(
			"/",
			uc.Update(),
		)
		apiRouter.Delete(
			"/",
			uc.Delete(),
		)
		apiRouter.With(uimwF).Post(
			"/uploaduserimage",
			uc.UploadUserImage(),
		)
		apiRouter.With(uimwT).Put(
			"/updateuserimage",
			uc.UpdateUserImage(),
		)
		apiRouter.With(uimwT).Delete(
			"/deleteuserimage",
			uc.DeleteUserImage(),
		)
	})
}

func EventRouter(r chi.Router, ev controllers.EventController, emw func(http.Handler) http.Handler) {
	imwT := middlewares.ImageMiddleware(true)
	imwF := middlewares.ImageMiddleware(false)
	isOwner := middlewares.IsOwnerMiddleware[domain.Event]()
	r.Route("/events", func(apiRouter chi.Router) {
		apiRouter.Post(
			"/",
			ev.Save(),
		)
		apiRouter.With(emw).Get(
			"/",
			ev.Find(),
		)
		apiRouter.Get(
			"/findlistby",
			ev.FindListBy(),
		)
		apiRouter.Get(
			"/findbyuser",
			ev.FindListByUser(),
		)
		apiRouter.With(emw, isOwner).Put(
			"/update",
			ev.Update(),
		)
		apiRouter.With(emw, isOwner).Delete(
			"/delete",
			ev.Delete(),
		)
		apiRouter.With(emw, imwF).Post(
			"/uploadimage",
			ev.UploadImage(),
		)
		apiRouter.With(emw, imwT).Put(
			"/updateimage",
			ev.UpdateImage(),
		)
		apiRouter.With(emw, imwT).Delete(
			"/deleteimage",
			ev.DeleteImage(),
		)
	})
}

func SubscriptionRouter(r chi.Router, ev controllers.SubscriptionController) {
	r.Route("/subscription", func(apiRouter chi.Router) {
		apiRouter.Post(
			"/",
			ev.Save(),
		)
		apiRouter.Delete(
			"/",
			ev.Delete(),
		)
		apiRouter.Get(
			"/",
			ev.FindUserSubs(),
		)
	})
}

func NotFoundJSON() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode("Resource Not Found")
		if err != nil {
			fmt.Printf("writing response: %s", err)
		}
	}
}

func PingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode("Ok")
		if err != nil {
			fmt.Printf("writing response: %s", err)
		}
	}
}
