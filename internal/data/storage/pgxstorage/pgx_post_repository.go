package pgxstorage

import (
	"context"
	"web_blog/internal/data/entity"
	"web_blog/internal/data/storage"

	"github.com/jackc/pgx"
)

type PgxPostRepository struct {
	Database *PgxDatabase
}

func (repository *PgxPostRepository) Create(ctx context.Context, tx *pgx.Tx, post *entity.Post) error {
	sql := `
		INSERT INTO posts (user_id, title, content) 
		VALUES ($1, $2, $3) 
		RETURNING id, verified, created_at, updated_at
	`
	return query(
		databasePayload[entity.Post]{
			conn: provideConn(tx, repository.Database.Connection),
			ctx:  ctx,
			sql:  sql,
			args: []any{post.UserID, post.Title, post.Content},
			scan: func(_ *entity.Post) []any {
				return []any{&post.ID, &post.Verified, &post.CreatedAt, &post.UpdatedAt}
			},
		},
	)
}

func (repository *PgxPostRepository) Find(ctx context.Context, tx *pgx.Tx, id int64) (*entity.Post, error) {
	sql := `
		SELECT * FROM posts WHERE id = $1
	`
	return queryOne(
		databasePayload[entity.Post]{
			conn: provideConn(tx, repository.Database.Connection),
			ctx:  ctx,
			sql:  sql,
			args: []any{id},
			scan: func(post *entity.Post) []any {
				return []any{
					&post.ID,
					&post.UserID,
					&post.Title,
					&post.Content,
					&post.Verified,
					&post.CreatedAt,
					&post.UpdatedAt,
				}
			},
		},
	)
}

func (repository *PgxPostRepository) FindAllByUserID(ctx context.Context, tx *pgx.Tx, filter storage.FilterQuery, id int64) ([]*entity.Post, error) {
	sql := `
		SELECT posts.* FROM posts
		INNER JOIN users ON posts.user_id = users.id
		WHERE users.id = $1
		LIMIT $2
		OFFSET $3
	`
	return queryAll(
		databasePayload[entity.Post]{
			conn: provideConn(tx, repository.Database.Connection),
			ctx:  ctx,
			sql:  sql,
			args: []any{id, filter.Limit, filter.Offset},
			scan: func(post *entity.Post) []any {
				return []any{
					&post.ID,
					&post.UserID,
					&post.Title,
					&post.Content,
					&post.Verified,
					&post.CreatedAt,
					&post.UpdatedAt,
				}
			},
		},
	)
}

func (repository *PgxPostRepository) FindAll(ctx context.Context, tx *pgx.Tx, filter storage.FilterQuery) ([]*entity.Post, error) {
	sql := `
		SELECT * FROM posts
		LIMIT $1
		OFFSET $2
	`
	return queryAll(
		databasePayload[entity.Post]{
			conn: provideConn(tx, repository.Database.Connection),
			ctx:  ctx,
			sql:  sql,
			args: []any{filter.Limit, filter.Offset},
			scan: func(post *entity.Post) []any {
				return []any{
					&post.ID,
					&post.UserID,
					&post.Title,
					&post.Content,
					&post.Verified,
					&post.CreatedAt,
					&post.UpdatedAt,
				}
			},
		},
	)
}

func (repository *PgxPostRepository) Update(ctx context.Context, tx *pgx.Tx, post *entity.Post) error {
	sql := `
		UPDATE posts 
		SET title=$1, content=$2, updated_at=NOW()
		WHERE id = $3
		RETURNING updated_at
		`
	return query(
		databasePayload[entity.Post]{
			conn: provideConn(tx, repository.Database.Connection),
			ctx:  ctx,
			sql:  sql,
			args: []any{post.Title, post.Content, post.ID},
			scan: func(_ *entity.Post) []any {
				return []any{&post.UpdatedAt}
			},
		},
	)
}

func (repository *PgxPostRepository) Delete(ctx context.Context, tx *pgx.Tx, id int64) error {
	sql := `
		DELETE FROM posts WHERE id = $1
	`
	return execute(
		databasePayload[entity.Post]{
			conn: provideConn(tx, repository.Database.Connection),
			ctx:  ctx,
			sql:  sql,
			args: []any{id},
			scan: nil,
		},
	)
}
