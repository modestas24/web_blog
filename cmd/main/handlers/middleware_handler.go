package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"web_blog/internal/auth"
	"web_blog/internal/data/entity"
	"web_blog/internal/data/storage"
)

type MiddlewareHandler struct {
	Storage       *storage.Storage
	Authenticator *auth.StatefulAuthenticator
	ErrorHandler  IErrorHandler
}

func (handler *MiddlewareHandler) StatefulAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user *entity.User
		var err error

		header := r.Header.Get("Authorization")
		if header == "" {
			handler.ErrorHandler.UnauthorizedError(w, r, errors.New("authorization header is missing"))
			return
		}

		values := strings.Split(header, " ")
		if len(values) < 2 {
			handler.ErrorHandler.UnauthorizedError(w, r, errors.New("authorization header is formated incorrectly"))
			return
		}

		ctx := r.Context()
		token := values[1]
		if user, err = handler.Authenticator.Validate(ctx, handler.Storage.Sessions, token); err != nil {
			handler.ErrorHandler.UnauthorizedError(w, r, errors.New("unauthorized"))
			return
		}

		ctx = context.WithValue(ctx, UserCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (handler *MiddlewareHandler) RoleMiddleware(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var user *entity.User
			var requiredRole *entity.Role
			var err error

			user = findUserFromCtx(r)
			if requiredRole, err = handler.Storage.Roles.FindByName(r.Context(), nil, role); err != nil {
				handler.ErrorHandler.InternalServerError(w, r, err)
				return
			}

			if user.Role.Level < requiredRole.Level {
				handler.ErrorHandler.ForbiddenError(w, r, errors.New("forbidden"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
