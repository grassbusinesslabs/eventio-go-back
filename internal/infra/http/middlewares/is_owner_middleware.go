package middlewares

import (
	"errors"
	"net/http"

	"github.com/grassbusinesslabs/eventio-go-back/internal/domain"
	"github.com/grassbusinesslabs/eventio-go-back/internal/infra/http/controllers"
)

type Userable interface {
	GetUserId() uint64
}

func IsOwnerMiddleware[domainType Userable]() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			user := ctx.Value(controllers.UserKey).(domain.User)
			obj := controllers.GetPathValFromCtx[domainType](ctx)

			if obj.GetUserId() != user.Id {
				err := errors.New("You have no access to this object!")
				controllers.Forbidden(w, err)
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}
