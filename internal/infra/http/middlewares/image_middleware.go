package middlewares

import (
	"errors"
	"net/http"

	"github.com/grassbusinesslabs/eventio-go-back/internal/domain"
	"github.com/grassbusinesslabs/eventio-go-back/internal/infra/http/controllers"
)

func ImageMiddleware(b bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			user := ctx.Value(controllers.UserKey).(domain.User)
			event := ctx.Value(controllers.EventKey).(domain.Event)
			if event.User_Id != user.Id {
				err := errors.New("You have no access to this image!")
				controllers.Forbidden(w, err)
				return
			}

			if b {
				if event.Image == "" {
					err := errors.New("Image doesn't exist!")
					controllers.Forbidden(w, err)
					return
				}
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}
