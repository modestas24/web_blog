package pgxstorage

import (
	"context"
	"web_blog/internal/data/entity"
	"web_blog/internal/data/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
)

type PgxUserRepository struct {
	Database *PgxDatabase
}

func (repository *PgxUserRepository) CreateWithVerification(
	ctx context.Context,
	tx *pgx.Tx,
	vrepository storage.IVerificationRepository,
	user *entity.User,
) error {
	return withTx(repository.Database.Connection, func(tx *pgx.Tx) error {
		var verification *entity.Verification
		var err error

		if err = repository.Create(ctx, tx, user); err != nil {
			return err
		}

		verification = &entity.Verification{UserID: user.ID}

		if err = vrepository.Create(ctx, tx, verification); err != nil {
			return err
		}

		return nil
	})
}

func (repository *PgxUserRepository) Verify(
	ctx context.Context,
	tx *pgx.Tx,
	vrepository storage.IVerificationRepository,
	id uuid.UUID,
	user *entity.User,
) error {
	return withTx(repository.Database.Connection, func(tx *pgx.Tx) error {
		var err error

		if err = repository.VerifyActive(ctx, tx, id, user); err != nil {
			return err
		}

		if err = vrepository.DeleteAllByUserID(ctx, tx, user.ID); err != nil {
			return err
		}

		return nil
	})
}

func (repository *PgxUserRepository) VerifyActive(ctx context.Context, tx *pgx.Tx, id uuid.UUID, user *entity.User) error {
	sql := `
		UPDATE users SET verified = true, updated_at = NOW()
		FROM verifications 
		WHERE users.id = verifications.user_id
		AND users.verified = false
		AND verifications.expired_at > NOW() 
		AND verifications.id = $1
		RETURNING users.id, users.role_id, users.email, users.username, users.verified, users.created_at, users.updated_at 
	`
	return query(
		databasePayload[entity.User]{
			conn: provideConn(tx, repository.Database.Connection),
			ctx:  ctx,
			sql:  sql,
			args: []any{id},
			scan: func(_ *entity.User) []any {
				return []any{
					&user.ID,
					&user.RoleID,
					&user.Email,
					&user.Username,
					&user.Verified,
					&user.CreatedAt,
					&user.UpdatedAt,
				}
			},
		},
	)
}

func (repository *PgxUserRepository) Create(ctx context.Context, tx *pgx.Tx, user *entity.User) error {
	sql := `
		INSERT INTO users (role_id, email, username, password) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, verified, created_at, updated_at
	`
	return query(
		databasePayload[entity.User]{
			conn: provideConn(tx, repository.Database.Connection),
			ctx:  ctx,
			sql:  sql,
			args: []any{user.RoleID, user.Email, user.Username, user.Password.Hash},
			scan: func(_ *entity.User) []any {
				return []any{
					&user.ID,
					&user.Verified,
					&user.CreatedAt,
					&user.UpdatedAt,
				}
			},
		},
	)
}

func (repository *PgxUserRepository) Find(ctx context.Context, tx *pgx.Tx, id int64) (*entity.User, error) {
	sql := `
		SELECT * FROM users WHERE id = $1
	`
	return queryOne(
		databasePayload[entity.User]{
			conn: provideConn(tx, repository.Database.Connection),
			ctx:  ctx,
			sql:  sql,
			args: []any{id},
			scan: func(user *entity.User) []any {
				return []any{
					&user.ID,
					&user.RoleID,
					&user.Email,
					&user.Username,
					&user.Password.Hash,
					&user.Verified,
					&user.CreatedAt,
					&user.UpdatedAt,
				}
			},
		},
	)
}

func (repository *PgxUserRepository) FindByEmail(ctx context.Context, tx *pgx.Tx, email string) (*entity.User, error) {
	sql := `
		SELECT * FROM users WHERE users.email = $1
	`
	return queryOne(
		databasePayload[entity.User]{
			conn: provideConn(tx, repository.Database.Connection),
			ctx:  ctx,
			sql:  sql,
			args: []any{email},
			scan: func(user *entity.User) []any {
				return []any{
					&user.ID,
					&user.RoleID,
					&user.Email,
					&user.Username,
					&user.Password.Hash,
					&user.Verified,
					&user.CreatedAt,
					&user.UpdatedAt,
				}
			},
		},
	)
}

func (repository *PgxUserRepository) FindAll(ctx context.Context, tx *pgx.Tx, filter storage.FilterQuery) ([]*entity.User, error) {
	sql := `
		SELECT * FROM users
		LIMIT $1
		OFFSET $2
	`
	return queryAll(
		databasePayload[entity.User]{
			conn: provideConn(tx, repository.Database.Connection),
			ctx:  ctx,
			sql:  sql,
			args: []any{filter.Limit, filter.Offset},
			scan: func(user *entity.User) []any {
				return []any{
					&user.ID,
					&user.RoleID,
					&user.Email,
					&user.Username,
					&user.Password.Hash,
					&user.Verified,
					&user.CreatedAt,
					&user.UpdatedAt,
				}
			},
		},
	)
}

func (repository *PgxUserRepository) Update(ctx context.Context, tx *pgx.Tx, user *entity.User) error {
	sql := `
		UPDATE users 
		SET role_id = $1, email = $2, username = $3, password = $4, verified = $5 updated_at = NOW()
		WHERE id = $6
		RETURNING updated_at
	`
	return query(
		databasePayload[entity.User]{
			conn: provideConn(tx, repository.Database.Connection),
			ctx:  ctx,
			sql:  sql,
			args: []any{user.RoleID, user.Email, user.Username, user.Password, user.Verified, user.ID},
			scan: func(_ *entity.User) []any {
				return []any{&user.UpdatedAt}
			},
		},
	)
}

func (repository *PgxUserRepository) Delete(ctx context.Context, tx *pgx.Tx, id int64) error {
	sql := `
		DELETE FROM users WHERE id = $1 
	`
	return execute(
		databasePayload[entity.User]{
			conn: provideConn(tx, repository.Database.Connection),
			ctx:  ctx,
			sql:  sql,
			args: []any{id},
			scan: nil,
		},
	)
}
