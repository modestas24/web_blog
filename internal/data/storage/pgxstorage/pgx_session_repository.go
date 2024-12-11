package pgxstorage

import (
	"context"
	"web_blog/internal/data/entity"
	"web_blog/internal/data/storage"

	"github.com/jackc/pgx"
)

type PgxSessionRepository struct {
	Database *PgxDatabase
}

func (repository *PgxSessionRepository) Create(ctx context.Context, tx *pgx.Tx, session *entity.Session) error {
	sql := `
		INSERT INTO sessions (id, user_id, expired_at) 
		VALUES ($1, $2, $3) 
	`
	return execute(
		databasePayload[entity.Session]{
			conn: provideConn(tx, repository.Database.Connection),
			ctx:  ctx,
			sql:  sql,
			args: []any{session.ID, session.UserID, session.ExpiredAt},
			scan: nil,
		},
	)
}

func (repository *PgxSessionRepository) Find(ctx context.Context, tx *pgx.Tx, id string) (*entity.Session, error) {
	sql := `
		SELECT * FROM sessions WHERE id = $1
	`
	return queryOne(
		databasePayload[entity.Session]{
			conn: provideConn(tx, repository.Database.Connection),
			ctx:  ctx,
			sql:  sql,
			args: []any{id},
			scan: func(session *entity.Session) []any {
				return []any{&session.ID, &session.UserID, &session.ExpiredAt}
			},
		},
	)
}

func (repository *PgxSessionRepository) FindWithUser(
	ctx context.Context,
	tx *pgx.Tx,
	id string,
) (*entity.Session, *entity.User, error) {
	type sessionWithUserPayload struct {
		session entity.Session
		user    entity.User
	}

	sql := `
		SELECT * FROM sessions
		INNER JOIN users ON sessions.user_id = users.id
		INNER JOIN roles ON users.role_id = roles.id
		WHERE sessions.id = $1
	`

	payload, err := queryOne(
		databasePayload[sessionWithUserPayload]{
			conn: provideConn(tx, repository.Database.Connection),
			ctx:  ctx,
			sql:  sql,
			args: []any{id},
			scan: func(payload *sessionWithUserPayload) []any {
				return []any{
					&payload.session.ID,
					&payload.session.UserID,
					&payload.session.ExpiredAt,

					&payload.user.ID,
					&payload.user.RoleID,
					&payload.user.Email,
					&payload.user.Username,
					&payload.user.Password.Hash,
					&payload.user.Verified,
					&payload.user.CreatedAt,
					&payload.user.UpdatedAt,

					&payload.user.Role.ID,
					&payload.user.Role.Level,
					&payload.user.Role.Name,
					&payload.user.Role.Description,
				}
			},
		},
	)

	return &payload.session, &payload.user, err
}

func (repository *PgxSessionRepository) FindAll(ctx context.Context, tx *pgx.Tx, filter storage.FilterQuery) ([]*entity.Session, error) {
	sql := `
		SELECT * FROM sessions
		LIMIT $1
		OFFSET $2
	`
	return queryAll(
		databasePayload[entity.Session]{
			conn: provideConn(tx, repository.Database.Connection),
			ctx:  ctx,
			sql:  sql,
			args: []any{filter.Limit, filter.Offset},
			scan: func(session *entity.Session) []any {
				return []any{&session.ID, &session.UserID, &session.ExpiredAt}
			},
		},
	)
}

func (repository *PgxSessionRepository) Update(ctx context.Context, tx *pgx.Tx, session *entity.Session) error {
	sql := `
		UPDATE sessions 
		SET id=$1, expired_at=$2
		WHERE id = $3
		RETURNING id, user_id, expired_at
		`
	return query(
		databasePayload[entity.Session]{
			conn: provideConn(tx, repository.Database.Connection),
			ctx:  ctx,
			sql:  sql,
			args: []any{session.ID, session.ExpiredAt, session.UserID},
			scan: func(_ *entity.Session) []any {
				return []any{&session.ID, &session.UserID, &session.ExpiredAt}
			},
		},
	)
}

func (repository *PgxSessionRepository) Delete(ctx context.Context, tx *pgx.Tx, id string) error {
	sql := `
		DELETE FROM sessions WHERE id = $1
	`
	return execute(
		databasePayload[entity.Session]{
			conn: provideConn(tx, repository.Database.Connection),
			ctx:  ctx,
			sql:  sql,
			args: []any{id},
			scan: nil,
		},
	)
}
