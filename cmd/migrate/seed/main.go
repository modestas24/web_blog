package main

import (
	"context"

	"web_blog/internal/data/storage"
	"web_blog/internal/data/storage/pgxstorage"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	godotenv.Load()
	ctx := context.TODO()
	logger := zap.Must(zap.NewProduction())

	var database pgxstorage.PgxDatabase
	if err := database.Open(ctx, nil); err != nil {
		logger.Fatal(err.Error())
		return
	}
	defer database.Close(ctx)

	store := &storage.Storage{
		Database:      &database,
		Users:         &pgxstorage.PgxUserRepository{Database: &database},
		Posts:         &pgxstorage.PgxPostRepository{Database: &database},
		Comments:      &pgxstorage.PgxCommentRepository{Database: &database},
		Verifications: &pgxstorage.PgxVerificationRepository{Database: &database},
	}

	logger.Info("seed has started")
	if err := storage.Seed(store); err != nil {
		logger.Warn("seeding error occured", zap.Error(err))
	}
	logger.Info("seed has ended")
}
