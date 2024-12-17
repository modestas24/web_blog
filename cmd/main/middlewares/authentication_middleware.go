package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"web_blog/cmd/main/utils"
	"web_blog/internal/data/entity"
)

type userKey string

const UserCtx userKey = "user"

func (middleware *Middleware) StatefulAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user *entity.User
		var err error

		header := r.Header.Get("Authorization")
		if header == "" {
			utils.UnauthorizedResponse(w, r, errors.New("authorization header is missing"))
			return
		}

		values := strings.Split(header, " ")
		if len(values) < 2 {
			utils.UnauthorizedResponse(w, r, errors.New("authorization header is formated incorrectly"))
			return
		}

		ctx := r.Context()
		token := values[1]
		if user, err = middleware.Authenticator.Validate(ctx, middleware.Storage.Sessions, token); err != nil {
			utils.UnauthorizedResponse(w, r, errors.New("unauthorized"))
			return
		}

		ctx = context.WithValue(ctx, UserCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func FindUserFromContext(r *http.Request) *entity.User {
	user, _ := r.Context().Value(UserCtx).(*entity.User)
	return user
}
