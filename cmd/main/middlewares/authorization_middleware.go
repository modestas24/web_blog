package middlewares

import (
	"errors"
	"net/http"
	"web_blog/cmd/main/utils"
	"web_blog/internal/data/entity"
)

func (middleware *Middleware) Authorization(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var user *entity.User
			var requiredRole *entity.Role
			var err error

			user = FindUserFromContext(r)
			if requiredRole, err = middleware.Storage.Roles.FindByName(r.Context(), nil, role); err != nil {
				utils.InternalServerErrorResponse(w, r, err)
				return
			}

			if user.Role.Level < requiredRole.Level {
				utils.ForbiddenResponse(w, r, errors.New("forbidden"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
