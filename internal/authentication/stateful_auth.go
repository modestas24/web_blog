package authentication

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"time"
	"web_blog/internal/data/entity"
	"web_blog/internal/data/storage"
)

type StatefulAuthenticator struct{}

func (authenticator *StatefulAuthenticator) Generate() (string, error) {
	var encoding *base32.Encoding
	var bytes []byte
	var err error

	encoding = base32.StdEncoding.WithPadding(base32.NoPadding)
	bytes = make([]byte, byteSize)
	if _, err = rand.Read(bytes); err != nil {
		return "", err
	}

	return encoding.EncodeToString(bytes), err
}

func (authenticator *StatefulAuthenticator) Create(
	ctx context.Context,
	repository storage.ISessionRepository,
	token string,
	id int64,
) error {
	var session *entity.Session
	var err error

	hash := sha256.Sum256([]byte(token))
	token = hex.EncodeToString(hash[:])

	session = &entity.Session{
		ID:        token,
		UserID:    id,
		ExpiredAt: time.Now().Add(expireDuration),
	}

	if err = repository.Create(ctx, nil, session); err != nil {
		return err
	}

	return nil
}

func (authenticator *StatefulAuthenticator) Validate(
	ctx context.Context,
	repository storage.ISessionRepository,
	token string,
) (*entity.User, error) {
	var user *entity.User
	var session *entity.Session
	var err error

	hash := sha256.Sum256([]byte(token))
	token = hex.EncodeToString(hash[:])

	if session, user, err = repository.FindWithUser(ctx, nil, token); err != nil {
		return nil, err
	}

	if time.Now().Compare(session.ExpiredAt) > 0 {
		if err = repository.Delete(ctx, nil, token); err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (authenticator *StatefulAuthenticator) Invalidate(
	ctx context.Context,
	repository storage.ISessionRepository,
	token string,
) error {
	var err error

	hash := sha256.Sum256([]byte(token))
	token = hex.EncodeToString(hash[:])

	if err = repository.Delete(ctx, nil, token); err != nil {
		return err
	}

	return nil
}
