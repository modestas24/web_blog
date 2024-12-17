package pgxstorage

import (
	"context"
	"time"
	"web_blog/internal/env"

	// Postgress database driver.
	"github.com/jackc/pgx"
)

// Postgress database.
type PgxDatabase struct {
	Connection *pgx.Conn
	Config     *pgx.ConnConfig
}

func (database *PgxDatabase) Open(ctx context.Context, config any) error {
	conf, ok := config.(pgx.ConnConfig)
	if !ok {
		conf = pgx.ConnConfig{
			Host:     env.GetString("DB_HOST", ""),
			Port:     uint16(env.GetInt("DB_PORT", 5432)),
			Database: env.GetString("DB_NAME", ""),
			User:     env.GetString("DB_USER", ""),
			Password: env.GetString("DB_PASSWORD", ""),
		}
	}

	database.Config = &conf
	conn, err := pgx.Connect(conf)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*500)
	defer cancel()
	if err := conn.Ping(ctx); err != nil {
		return err
	}

	database.Connection = conn
	return nil
}

func (database *PgxDatabase) Close(ctx context.Context) error {
	return database.Connection.Close()
}
