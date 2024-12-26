package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/grassbusinesslabs/eventio-go-back/internal/domain"
	"github.com/grassbusinesslabs/eventio-go-back/internal/infra/http/controllers"
)

func AvatarMiddleware(b bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			user := ctx.Value(controllers.UserKey).(domain.User)

			filePath := fmt.Sprintf("file_storage/user_image/%d.png", user.Id)
			_, err := os.Stat(filePath)

			if b {
				if os.IsNotExist(err) {
					err := errors.New("Image doesn't exist!")
					controllers.Forbidden(w, err)
					return
				}
			} else {
				if err == nil {
					err := errors.New("Image already exists!")
					controllers.Forbidden(w, err)
					return
				}
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
