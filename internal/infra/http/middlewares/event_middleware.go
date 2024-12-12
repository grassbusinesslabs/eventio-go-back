package middlewares

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/grassbusinesslabs/eventio-go-back/internal/app"
	"github.com/grassbusinesslabs/eventio-go-back/internal/infra/http/controllers"
)

func EventMiddleware(es app.EventService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			eventId, err := strconv.ParseUint(chi.URLParam(r, "eventid"), 10, 64)
			if err != nil {
				log.Printf("EventMiddleware -> strconv.ParseUint: %s", err)
				controllers.BadRequest(w, err)
				return
			}

			event, err := es.Find(eventId)
			if err != nil {
				log.Printf("EventMiddleware ->  es.Find: %s", err)
				controllers.InternalServerError(w, err)
				return
			}

			ctx = context.WithValue(ctx, controllers.EventKey, event)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}
