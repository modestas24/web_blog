package api

import (
	"net/http"
	"time"
	"web_blog/cmd/main/handlers"
	"web_blog/docs"
	"web_blog/internal/auth"
	"web_blog/internal/data/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
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
	Handlers      handlers.Handlers
	Storage       storage.Storage
	Logger        *zap.Logger
	Authenticator *auth.StatefulAuthenticator
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

func (application *Application) Mount() http.Handler {
	handlers := application.Handlers
	authMiddleware := handlers.Middleware.StatefulAuthMiddleware
	roleMiddleware := handlers.Middleware.RoleMiddleware
	postMiddleware := handlers.Post.PostContextMiddleware

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", handlers.Health.CheckHealthHandler)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(application.Config.SwaggerConfig.DocsURL)))

		// Post handlers
		r.Group(func(r chi.Router) {
			r.Get("/posts", handlers.Post.FindAllPostHandler)
			r.Get("/users/{id}/posts", handlers.Post.FindAllPostByUserIDHandler)
			r.With(postMiddleware).Get("/posts/{id}", handlers.Post.FindPostHandler)

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware)

				r.With(roleMiddleware("user")).Post("/posts", handlers.Post.CreatePostHandler)
				r.With(roleMiddleware("moderator"), postMiddleware).Patch("/posts/{id}", handlers.Post.UpdatePostHandler)
				r.With(roleMiddleware("moderator"), postMiddleware).Delete("/posts/{id}", handlers.Post.DeletePostHandler)
			})

		})

		// Comment handlers
		r.Group(func(r chi.Router) {
			r.Get("/posts/{id}/comments", handlers.Comment.FindAllByPostIDCommentHandler)

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware)
				r.With(roleMiddleware("user")).Post("/posts/{id}/comments", handlers.Comment.CreateCommentHandler)
				r.With(roleMiddleware("moderator")).Get("/posts/comments", handlers.Comment.FindAllCommentHandler)
				r.With(roleMiddleware("moderator")).Delete("/posts/comments/{id}", handlers.Comment.DeleteCommentHandler)
			})
		})

		// User handlers
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware, roleMiddleware("admin"))
			r.Get("/users", handlers.User.FindAllUsersHandler)
		})

		// Authentication handlers
		r.Group(func(r chi.Router) {
			r.Post("/authentication/register", handlers.Auth.RegisterUserHandler)
			r.Post("/authentication/verify", handlers.Auth.VerifyUserHandler)
			r.Post("/authentication/login", handlers.Auth.LoginUserHandler)
			r.Delete("/authentication/logout", handlers.Auth.LogoutUserHandler)
		})
	})

	return r
}

func (application *Application) Serve() error {
	docs.SwaggerInfo.Title = Title
	docs.SwaggerInfo.Description = Description
	docs.SwaggerInfo.BasePath = BasePath
	docs.SwaggerInfo.Version = Version
	docs.SwaggerInfo.Host = application.Config.Url

	srv := &http.Server{
		Addr:         application.Config.Address,
		Handler:      application.Mount(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	application.Logger.Info("Server has started", zap.String("address", application.Config.Address))
	return srv.ListenAndServe()
}
