package pgxstorage

import (
	"context"
	"web_blog/internal/data/entity"
	"web_blog/internal/data/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
)

type PgxVerificationRepository struct {
	Database *PgxDatabase
}

func (repository *PgxVerificationRepository) Create(ctx context.Context, tx *pgx.Tx, verification *entity.Verification) error {
	sql := `
		INSERT INTO verifications (user_id) 
		VALUES ($1) 
		RETURNING id, expired_at
	`
	return query(
		databasePayload[entity.Verification]{
			conn: provideConn(tx, repository.Database.Connection),
			ctx:  ctx,
			sql:  sql,
			args: []any{verification.UserID},
			scan: func(_ *entity.Verification) []any {
				return []any{&verification.UUID, &verification.ExpiredAt}
			},
		},
	)
}

func (repository *PgxVerificationRepository) Find(ctx context.Context, tx *pgx.Tx, id uuid.UUID) (*entity.Verification, error) {
	sql := `
		SELECT * FROM verifications WHERE id = $1
	`
	return queryOne(
		databasePayload[entity.Verification]{
			conn: provideConn(tx, repository.Database.Connection),
			ctx:  ctx,
			sql:  sql,
			args: []any{id},
			scan: func(verifiction *entity.Verification) []any {
				return []any{&verifiction.UUID, &verifiction.UserID, &verifiction.ExpiredAt}
			},
		},
	)
}

func (repository *PgxVerificationRepository) FindAllByUserID(
	ctx context.Context,
	tx *pgx.Tx,
	filter storage.FilterQuery,
	id int64,
) ([]*entity.Verification, error) {
	sql := `
		SELECT * FROM verifications WHERE id = $1
		LIMIT $2
		OFFSET $3
	`
	return queryAll(
		databasePayload[entity.Verification]{
			conn: provideConn(tx, repository.Database.Connection),
			ctx:  ctx,
			sql:  sql,
			args: []any{id, filter.Limit, filter.Offset},
			scan: func(verification *entity.Verification) []any {
				return []any{&verification.UUID, &verification.UserID, &verification.ExpiredAt}
			},
		},
	)
}

func (repository *PgxVerificationRepository) FindAll(ctx context.Context, tx *pgx.Tx, filter storage.FilterQuery) ([]*entity.Verification, error) {
	sql := `
		SELECT * FROM verifications
		LIMIT $1
		OFFSET $2
	`
	return queryAll(
		databasePayload[entity.Verification]{
			conn: provideConn(tx, repository.Database.Connection),
			ctx:  ctx,
			sql:  sql,
			args: []any{filter.Limit, filter.Offset},
			scan: func(verification *entity.Verification) []any {
				return []any{&verification.UUID, &verification.UserID, &verification.ExpiredAt}
			},
		},
	)
}

func (repository *PgxVerificationRepository) Update(ctx context.Context, tx *pgx.Tx, user_ver *entity.Verification) error {
	return nil
}

func (repository *PgxVerificationRepository) Delete(ctx context.Context, tx *pgx.Tx, id uuid.UUID) error {
	sql := `
		DELETE FROM verifications WHERE id = $1 
	`
	return execute(
		databasePayload[entity.Verification]{
			conn: provideConn(tx, repository.Database.Connection),
			ctx:  ctx,
			sql:  sql,
			args: []any{id},
			scan: nil,
		},
	)
}

func (repository *PgxVerificationRepository) DeleteAllByUserID(ctx context.Context, tx *pgx.Tx, id int64) error {
	sql := `
		DELETE FROM verifications WHERE user_id = $1 
	`
	return execute(
		databasePayload[entity.Verification]{
			conn: provideConn(tx, repository.Database.Connection),
			ctx:  ctx,
			sql:  sql,
			args: []any{id},
			scan: nil,
		},
	)
}
