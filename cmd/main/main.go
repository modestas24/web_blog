package main

import (
	"context"
	"fmt"

	"web_blog/cmd/main/api"
	handlers "web_blog/cmd/main/handlers"
	"web_blog/internal/auth"
	"web_blog/internal/data/storage"
	pgxstorage "web_blog/internal/data/storage/pgxstorage"
	"web_blog/internal/env"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

//	@title			Golang Web Blog API
//	@description	Blog API written in Golang for university module.

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description				User token required for authorization
func main() {
	godotenv.Load()
	ctx := context.TODO()

	// Logger
	logger := zap.Must(zap.NewProduction())

	// Database
	database := &pgxstorage.PgxDatabase{}
	if err := database.Open(ctx, nil); err != nil {
		logger.Fatal(err.Error())
		return
	}
	defer database.Close(ctx)

	// Authenticator
	authenticator := &auth.StatefulAuthenticator{}

	storage := storage.Storage{
		Database:      database,
		Users:         &pgxstorage.PgxUserRepository{Database: database},
		Posts:         &pgxstorage.PgxPostRepository{Database: database},
		Comments:      &pgxstorage.PgxCommentRepository{Database: database},
		Verifications: &pgxstorage.PgxVerificationRepository{Database: database},
		Sessions:      &pgxstorage.PgxSessionRepository{Database: database},
		Roles:         &pgxstorage.PgxRoleRepository{Database: database},
	}

	// Error handler
	errorHandler := &handlers.ErrorHandler{
		Logger: logger,
	}

	// Handlers
	handlers := handlers.Handlers{
		Errors: errorHandler,
		Middleware: &handlers.MiddlewareHandler{
			Storage:       &storage,
			Authenticator: authenticator,
			ErrorHandler:  errorHandler,
		},
		Health: &handlers.HealthHandler{
			HealthEnvelope: handlers.HealthEnvelope{
				Title:       api.Title,
				Description: api.Description,
				Version:     api.Version,
				Url:         env.GetString("URL", "localhost:8080"),
			},
		},
		Auth: &handlers.AuthHandler{
			Storage:         &storage,
			Authentificator: authenticator,
			ErrorHandler:    errorHandler,
		},
		User: &handlers.UserHandler{
			Storage:      &storage,
			ErrorHandler: errorHandler,
		},
		Post: &handlers.PostHandler{
			Storage:      &storage,
			ErrorHandler: errorHandler,
		},
		Comment: &handlers.CommentHandler{
			Storage:      &storage,
			ErrorHandler: errorHandler,
		},
	}

	// Application config
	address := env.GetString("ADDR", ":8080")
	url := env.GetString("URL", "localhost:8080")
	config := api.Config{
		Address: address,
		Url:     url,
		Storage: database.Config,
		SwaggerConfig: api.SwaggerConfig{
			DocsURL: fmt.Sprintf("http://%s%s/swagger/doc.json", address, api.BasePath),
		},
	}

	application := &api.Application{
		Config:        config,
		Handlers:      handlers,
		Storage:       storage,
		Logger:        logger,
		Authenticator: authenticator,
	}

	logger.Fatal(application.Serve().Error())
}
