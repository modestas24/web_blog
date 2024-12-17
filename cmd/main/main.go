package main

import (
	"context"
	"fmt"
	"web_blog/cmd/main/api"
	"web_blog/cmd/main/middlewares"
	"web_blog/cmd/main/services"
	"web_blog/internal/authentication"
	"web_blog/internal/data/storage"
	"web_blog/internal/data/storage/pgxstorage"
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
	var err error

	// Logger
	Logger := zap.Must(zap.NewProduction())

	if err = godotenv.Load(); err != nil {
		Logger.Fatal("dotenv error", zap.Error(err))
	}

	url := env.GetString("URL", "localhost:8080")
	address := env.GetString("ADDR", "localhost:8080")
	healthEnvelope := services.HealthEnvelope{
		Title:       api.Title,
		Description: api.Description,
		Version:     api.Version,
		Url:         address,
	}

	// Database
	Database := &pgxstorage.PgxDatabase{}
	if err = Database.Open(context.Background(), nil); err != nil {
		Logger.Fatal("database error", zap.Error(err))
		return
	}
	defer Database.Close(context.Background())

	// Authenticator
	Authenticator := authentication.StatefulAuthenticator{}

	// Storage
	Storage := storage.Storage{
		Database:      Database,
		Users:         &pgxstorage.PgxUserRepository{Database: Database},
		Posts:         &pgxstorage.PgxPostRepository{Database: Database},
		Comments:      &pgxstorage.PgxCommentRepository{Database: Database},
		Verifications: &pgxstorage.PgxVerificationRepository{Database: Database},
		Sessions:      &pgxstorage.PgxSessionRepository{Database: Database},
		Roles:         &pgxstorage.PgxRoleRepository{Database: Database},
	}

	// Middlewares
	Middlewares := middlewares.Middleware{
		Storage:       &Storage,
		Authenticator: &Authenticator,
	}

	// Services
	Services := services.Services{
		Health:  &services.HealthService{HealthEnvelope: healthEnvelope},
		Auth:    &services.AuthService{Storage: &Storage, Authenticator: &Authenticator},
		User:    &services.UserService{Storage: &Storage},
		Post:    &services.PostService{Storage: &Storage},
		Comment: &services.CommentService{Storage: &Storage},
	}

	// Application config
	Config := api.Config{
		Address: address,
		Url:     url,
		Storage: Database.Config,
		SwaggerConfig: api.SwaggerConfig{
			DocsURL: fmt.Sprintf("http://%s%s/swagger/doc.json", address, api.BasePath),
		},
	}

	// Application
	Application := &api.Application{
		Config:        Config,
		Middlewares:   Middlewares,
		Services:      Services,
		Storage:       Storage,
		Logger:        Logger,
		Authenticator: Authenticator,
	}

	Logger.Fatal(Application.Serve().Error())
}
