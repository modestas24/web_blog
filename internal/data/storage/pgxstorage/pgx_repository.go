package pgxstorage

import (
	"context"
	"web_blog/internal/data/storage"

	"github.com/jackc/pgx"
)

type connection interface {
	ExecEx(context.Context, string, *pgx.QueryExOptions, ...interface{}) (pgx.CommandTag, error)
	QueryEx(context.Context, string, *pgx.QueryExOptions, ...interface{}) (*pgx.Rows, error)
	QueryRowEx(context.Context, string, *pgx.QueryExOptions, ...interface{}) *pgx.Row
}

type databasePayload[T any] struct {
	conn connection
	ctx  context.Context
	sql  string
	args []any
	scan func(*T) []any
}

func provideConn(conn *pgx.Tx, alt *pgx.Conn) connection {
	if conn != nil {
		return conn
	}

	return alt
}

func withTx(conn *pgx.Conn, fn func(*pgx.Tx) error) error {
	var tx *pgx.Tx
	var err error

	if tx, err = conn.Begin(); err != nil {
		return err
	}

	if err = fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func execute[T any](dp databasePayload[T]) error {
	var com pgx.CommandTag
	var err error

	ctx, cancel := context.WithTimeout(dp.ctx, storage.DatabaseQueryTimeout)
	defer cancel()

	if com, err = dp.conn.ExecEx(ctx, dp.sql, nil, dp.args...); err != nil {
		return err
	}

	if com.RowsAffected() <= 0 {
		return storage.ErrorNotFound
	}

	return nil
}

func query[T any](dp databasePayload[T]) error {
	ctx, cancel := context.WithTimeout(dp.ctx, storage.DatabaseQueryTimeout)
	defer cancel()

	return dp.conn.
		QueryRowEx(ctx, dp.sql, nil, dp.args...).
		Scan(dp.scan(nil)...)
}

func queryOne[T any](dp databasePayload[T]) (*T, error) {
	ctx, cancel := context.WithTimeout(dp.ctx, storage.DatabaseQueryTimeout)
	defer cancel()

	element := new(T)

	return element, dp.conn.
		QueryRowEx(ctx, dp.sql, nil, dp.args...).
		Scan(dp.scan(element)...)
}

func queryAll[T any](dp databasePayload[T]) ([]*T, error) {
	var rows *pgx.Rows
	var list []*T
	var err error

	ctx, cancel := context.WithTimeout(dp.ctx, storage.DatabaseQueryTimeout)
	defer cancel()

	if rows, err = dp.conn.QueryEx(ctx, dp.sql, nil, dp.args...); err != nil {
		return list, err
	}

	defer rows.Close()
	for rows.Next() {
		element := new(T)
		if err = rows.Scan(dp.scan(element)...); err != nil {
			return list, err
		}
		list = append(list, element)
	}

	return list, err
}
