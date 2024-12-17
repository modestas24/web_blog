package api

import (
	"net/http"
	"time"
	"web_blog/cmd/main/middlewares"
	"web_blog/cmd/main/services"
	"web_blog/docs"
	"web_blog/internal/authentication"
	"web_blog/internal/data/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

const (
	Title       string = "Golang Web Blog API"
	Description string = "Blog API written in Golang for university module."
	Version     string = "0.1"
	BasePath    string = "/v1"
)

type Application struct {
	Config        Config
	Middlewares   middlewares.Middleware
	Services      services.Services
	Storage       storage.Storage
	Authenticator authentication.StatefulAuthenticator
	Logger        *zap.Logger
}

type Config struct {
	Address       string
	Url           string
	Storage       any
	SwaggerConfig SwaggerConfig
}

type SwaggerConfig struct {
	Title       string
	Description string
	BasePath    string
	Version     string
	Host        string
	DocsURL     string
}

func (app *Application) Mount() http.Handler {
	Services := app.Services
	Middlewares := app.Middlewares

	StatefulAuthentication := Middlewares.StatefulAuthentication
	Authorization := Middlewares.Authorization
	PostContext := Middlewares.PostContext

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", Services.Health.CheckHealth)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(app.Config.SwaggerConfig.DocsURL)))

		// Post Services.
		r.Group(func(r chi.Router) {
			r.Get("/posts", Services.Post.FindAllPosts)
			r.Get("/users/{id}/posts", Services.Post.FindAllPostsByUserID)
			r.With(PostContext).
				Get("/posts/{id}", Services.Post.FindPost)

			// With Authentication.
			r.Group(func(r chi.Router) {
				r.Use(StatefulAuthentication)

				r.With(Authorization("user")).
					Post("/posts", Services.Post.CreatePost)
				r.With(Authorization("moderator"), PostContext).
					Patch("/posts/{id}", Services.Post.UpdatePost)
				r.With(Authorization("moderator"), PostContext).
					Delete("/posts/{id}", Services.Post.DeletePost)
			})

		})

		// Comment Services.
		r.Group(func(r chi.Router) {
			r.Get("/posts/{id}/comments", Services.Comment.FindAllCommentsByPostID)

			// With Authentication.
			r.Group(func(r chi.Router) {
				r.Use(StatefulAuthentication)
				r.With(Authorization("user")).
					Post("/posts/{id}/comments", Services.Comment.CreateComment)
				r.With(Authorization("moderator")).
					Get("/posts/comments", Services.Comment.FindAllComments)
				r.With(Authorization("moderator")).
					Delete("/posts/comments/{id}", Services.Comment.DeleteComment)
			})
		})

		// User Services.
		r.Group(func(r chi.Router) {
			r.Use(StatefulAuthentication, Authorization("admin"))
			r.Get("/users", Services.User.FindAllUsers)
		})

		// Authentication Services.
		r.Group(func(r chi.Router) {
			r.Post("/authentication/register", Services.Auth.RegisterUser)
			r.Post("/authentication/verify", Services.Auth.VerifyUser)
			r.Post("/authentication/login", Services.Auth.LoginUser)
			r.Delete("/authentication/logout", Services.Auth.LogoutUser)
		})
	})

	return r
}

func (app *Application) Serve() error {
	docs.SwaggerInfo.Title = Title
	docs.SwaggerInfo.Description = Description
	docs.SwaggerInfo.BasePath = BasePath
	docs.SwaggerInfo.Version = Version
	docs.SwaggerInfo.Host = app.Config.Url

	srv := &http.Server{
		Addr:         app.Config.Address,
		Handler:      app.Mount(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	app.Logger.Info("Server has started", zap.String("address", app.Config.Address))
	return srv.ListenAndServe()
}
