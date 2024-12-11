package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
)

type IErrorHandler interface {
	InternalServerError(http.ResponseWriter, *http.Request, error)
	BadRequestError(http.ResponseWriter, *http.Request, error)
	NotFoundError(http.ResponseWriter, *http.Request, error)
	ForbiddenError(http.ResponseWriter, *http.Request, error)
	UnauthorizedError(http.ResponseWriter, *http.Request, error)
	SwitchInternalServerError(http.ResponseWriter, *http.Request, error)
}

type IMiddlewareHandler interface {
	StatefulAuthMiddleware(http.Handler) http.Handler
	RoleMiddleware(string) func(http.Handler) http.Handler
}

type IHealthHandler interface {
	CheckHealthHandler(w http.ResponseWriter, _ *http.Request)
}

type IAuthHandler interface {
	RegisterUserHandler(http.ResponseWriter, *http.Request)
	VerifyUserHandler(http.ResponseWriter, *http.Request)
	LoginUserHandler(http.ResponseWriter, *http.Request)
	LogoutUserHandler(http.ResponseWriter, *http.Request)
}

type IUserHandler interface {
	FindAllUsersHandler(http.ResponseWriter, *http.Request)
}

type IPostHandler interface {
	CreatePostHandler(http.ResponseWriter, *http.Request)
	FindAllPostHandler(http.ResponseWriter, *http.Request)
	FindAllPostByUserIDHandler(http.ResponseWriter, *http.Request)
	FindPostHandler(http.ResponseWriter, *http.Request)
	UpdatePostHandler(http.ResponseWriter, *http.Request)
	DeletePostHandler(http.ResponseWriter, *http.Request)

	// Middlewares
	PostContextMiddleware(next http.Handler) http.Handler
}

type ICommentHandler interface {
	CreateCommentHandler(http.ResponseWriter, *http.Request)
	FindAllCommentHandler(http.ResponseWriter, *http.Request)
	FindAllByPostIDCommentHandler(http.ResponseWriter, *http.Request)
	DeleteCommentHandler(http.ResponseWriter, *http.Request)
}

var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

type Handlers struct {
	Errors     IErrorHandler
	Middleware IMiddlewareHandler
	Health     IHealthHandler
	Auth       IAuthHandler
	User       IUserHandler
	Post       IPostHandler
	Comment    ICommentHandler
}
