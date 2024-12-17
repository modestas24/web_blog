package middlewares

import (
	"web_blog/internal/authentication"
	"web_blog/internal/data/storage"
)

type Middleware struct {
	Storage       *storage.Storage
	Authenticator *authentication.StatefulAuthenticator
}
