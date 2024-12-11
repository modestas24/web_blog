package pgxstorage

import (
	"context"
	"web_blog/internal/data/entity"
	"web_blog/internal/data/storage"

	"github.com/jackc/pgx"
)

type PgxCommentRepository struct {
	Database *PgxDatabase
}

func (repository *PgxCommentRepository) Create(ctx context.Context, tx *pgx.Tx, comment *entity.Comment) error {
	sql := `
		INSERT INTO comments (user_id, post_id, content) 
		VALUES ($1, $2, $3) 
		RETURNING id, verified, created_at, updated_at
	`
	return query(
		databasePayload[entity.Comment]{
			conn: provideConn(tx, repository.Database.Connection),
			ctx:  ctx,
			sql:  sql,
			args: []any{comment.UserID, comment.PostID, comment.Content},
			scan: func(_ *entity.Comment) []any {
				return []any{&comment.ID, &comment.Verified, &comment.CreatedAt, &comment.UpdatedAt}
			},
		},
	)
}

func (repository *PgxCommentRepository) Find(ctx context.Context, tx *pgx.Tx, id int64) (*entity.Comment, error) {
	sql := `
		SELECT * FROM comments WHERE id = $1
	`
	return queryOne(
		databasePayload[entity.Comment]{
			conn: provideConn(tx, repository.Database.Connection),
			ctx:  ctx,
			sql:  sql,
			args: []any{id},
			scan: func(c *entity.Comment) []any {
				return []any{&c.ID, &c.UserID, &c.PostID, &c.Content, &c.Verified, &c.CreatedAt, &c.UpdatedAt}
			},
		},
	)
}

func (repository *PgxCommentRepository) FindAll(ctx context.Context, tx *pgx.Tx, filter storage.FilterQuery) ([]*entity.Comment, error) {
	sql := `
		SELECT * FROM comments
		LIMIT $1
		OFFSET $2
	`
	return queryAll(
		databasePayload[entity.Comment]{
			conn: provideConn(tx, repository.Database.Connection),
			ctx:  ctx,
			sql:  sql,
			args: []any{filter.Limit, filter.Offset},
			scan: func(c *entity.Comment) []any {
				return []any{&c.ID, &c.UserID, &c.PostID, &c.Content, &c.Verified, &c.CreatedAt, &c.UpdatedAt}
			},
		},
	)
}

func (repository *PgxCommentRepository) FindAllByUserID(ctx context.Context, tx *pgx.Tx, filter storage.FilterQuery, id int64) ([]*entity.Comment, error) {
	sql := `
		SELECT comments.* FROM comments
		INNER JOIN users ON comments.user_id = users.id
		WHERE users.id = $1
		LIMIT $2
		OFFSET $3
	`
	return queryAll(
		databasePayload[entity.Comment]{
			conn: provideConn(tx, repository.Database.Connection),
			ctx:  ctx,
			sql:  sql,
			args: []any{id, filter.Limit, filter.Offset},
			scan: func(c *entity.Comment) []any {
				return []any{&c.ID, &c.UserID, &c.PostID, &c.Content, &c.Verified, &c.CreatedAt, &c.UpdatedAt}
			},
		},
	)
}

func (repository *PgxCommentRepository) FindAllByPostID(ctx context.Context, tx *pgx.Tx, filter storage.FilterQuery, id int64) ([]*entity.Comment, error) {
	sql := `
		SELECT comments.* FROM comments
		INNER JOIN posts ON comments.post_id = posts.id 
		WHERE posts.id = $1
		LIMIT $2
		OFFSET $3
	`
	return queryAll(
		databasePayload[entity.Comment]{
			conn: provideConn(tx, repository.Database.Connection),
			ctx:  ctx,
			sql:  sql,
			args: []any{id, filter.Limit, filter.Offset},
			scan: func(c *entity.Comment) []any {
				return []any{&c.ID, &c.UserID, &c.PostID, &c.Content, &c.Verified, &c.CreatedAt, &c.UpdatedAt}
			},
		},
	)
}

func (repository *PgxCommentRepository) Update(ctx context.Context, tx *pgx.Tx, comment *entity.Comment) error {
	sql := `
		UPDATE comments 
		SET content = $1, updated_at = NOW() 
		WHERE id = $2
		RETURNING updated_at
	`
	return query(
		databasePayload[entity.Comment]{
			conn: provideConn(tx, repository.Database.Connection),
			ctx:  ctx,
			sql:  sql,
			args: []any{comment.Content, comment.ID},
			scan: func(comment *entity.Comment) []any {
				return []any{&comment.UpdatedAt}
			},
		},
	)
}

func (repository *PgxCommentRepository) Delete(ctx context.Context, tx *pgx.Tx, id int64) error {
	sql := `
		DELETE FROM comments WHERE id = $1
	`
	return execute(
		databasePayload[entity.Comment]{
			conn: provideConn(tx, repository.Database.Connection),
			ctx:  ctx,
			sql:  sql,
			args: []any{id},
			scan: nil,
		},
	)
}
