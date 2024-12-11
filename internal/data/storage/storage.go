package storage

import (
	"context"
	"errors"
	"time"
	"web_blog/internal/data/entity"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
)

var (
	ErrorNotFound        error         = errors.New("resource not found")
	ErrorDuplicate       error         = errors.New("resource already exists")
	DatabaseQueryTimeout time.Duration = time.Second * 3
)

type Database interface {
	Open(context.Context, any) error
	Close(context.Context) error
}

type IRepository[T any, ID any] interface {
	Create(context.Context, *pgx.Tx, *T) error
	Find(context.Context, *pgx.Tx, ID) (*T, error)
	FindAll(context.Context, *pgx.Tx, FilterQuery) ([]*T, error)
	Update(context.Context, *pgx.Tx, *T) error
	Delete(context.Context, *pgx.Tx, ID) error
}

type IUserRepository interface {
	IRepository[entity.User, int64]
	CreateWithVerification(context.Context, *pgx.Tx, IVerificationRepository, *entity.User) error
	Verify(context.Context, *pgx.Tx, IVerificationRepository, uuid.UUID, *entity.User) error
	FindByEmail(context.Context, *pgx.Tx, string) (*entity.User, error)
}

type IPostRepository interface {
	IRepository[entity.Post, int64]
	FindAllByUserID(context.Context, *pgx.Tx, FilterQuery, int64) ([]*entity.Post, error)
}

type ICommentRepository interface {
	IRepository[entity.Comment, int64]
	FindAllByUserID(context.Context, *pgx.Tx, FilterQuery, int64) ([]*entity.Comment, error)
	FindAllByPostID(context.Context, *pgx.Tx, FilterQuery, int64) ([]*entity.Comment, error)
}

type IVerificationRepository interface {
	IRepository[entity.Verification, uuid.UUID]
	FindAllByUserID(context.Context, *pgx.Tx, FilterQuery, int64) ([]*entity.Verification, error)
	DeleteAllByUserID(context.Context, *pgx.Tx, int64) error
}

type ISessionRepository interface {
	IRepository[entity.Session, string]
	FindWithUser(context.Context, *pgx.Tx, string) (*entity.Session, *entity.User, error)
}

type IRoleRepository interface {
	IRepository[entity.Role, int64]
	FindByName(context.Context, *pgx.Tx, string) (*entity.Role, error)
}

type Storage struct {
	Database      Database
	Users         IUserRepository
	Posts         IPostRepository
	Comments      ICommentRepository
	Verifications IVerificationRepository
	Sessions      ISessionRepository
	Roles         IRoleRepository
}
