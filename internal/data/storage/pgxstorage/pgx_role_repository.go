package pgxstorage

import (
	"context"
	"web_blog/internal/data/entity"
	"web_blog/internal/data/storage"

	"github.com/jackc/pgx"
)

type PgxRoleRepository struct {
	Database *PgxDatabase
}

func (repository *PgxRoleRepository) Create(ctx context.Context, tx *pgx.Tx, session *entity.Role) error {
	return nil
}

func (repository *PgxRoleRepository) Find(ctx context.Context, tx *pgx.Tx, id int64) (*entity.Role, error) {
	return nil, nil
}

func (repository *PgxRoleRepository) FindByName(ctx context.Context, tx *pgx.Tx, name string) (*entity.Role, error) {
	sql := `
		SELECT * FROM roles WHERE name = $1
	`
	return queryOne(
		databasePayload[entity.Role]{
			conn: provideConn(tx, repository.Database.Connection),
			ctx:  ctx,
			sql:  sql,
			args: []any{name},
			scan: func(role *entity.Role) []any {
				return []any{
					&role.ID,
					&role.Level,
					&role.Name,
					&role.Description,
				}
			},
		},
	)
}

func (repository *PgxRoleRepository) FindAll(ctx context.Context, tx *pgx.Tx, filter storage.FilterQuery) ([]*entity.Role, error) {
	return nil, nil
}

func (repository *PgxRoleRepository) Update(ctx context.Context, tx *pgx.Tx, session *entity.Role) error {
	return nil
}

func (repository *PgxRoleRepository) Delete(ctx context.Context, tx *pgx.Tx, id int64) error {
	return nil
}
