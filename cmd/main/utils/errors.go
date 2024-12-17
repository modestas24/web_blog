package utils

import (
	"errors"
	"net/http"
	"web_blog/internal/data/storage"

	"github.com/jackc/pgx"
	"go.uber.org/zap"
)

type ErrorService struct {
	Logger *zap.Logger
}

var logger = zap.Must(zap.NewProduction())

func writeResponse(w http.ResponseWriter, r *http.Request, status int, msg string, err error) {
	logger.Warn(msg, zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.Error(err))
	if err := WriteJsonError(w, r, status, err.Error()); err != nil {
		logger.Error(err.Error())
	}
}

func InternalServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	writeResponse(w, r, http.StatusInternalServerError, "internal server error", err)
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	writeResponse(w, r, http.StatusBadRequest, "bad request error", err)
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	writeResponse(w, r, http.StatusNotFound, "not found error", err)
}

func ForbiddenResponse(w http.ResponseWriter, r *http.Request, err error) {
	writeResponse(w, r, http.StatusForbidden, "forbidden error", err)
}

func UnauthorizedResponse(w http.ResponseWriter, r *http.Request, err error) {
	writeResponse(w, r, http.StatusUnauthorized, "unauthorized error", err)
}

func SwitchInternalServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, storage.ErrorNotFound):
		NotFoundResponse(w, r, err)
		return
	case errors.Is(err, pgx.ErrNoRows):
		NotFoundResponse(w, r, storage.ErrorNotFound)
		return
	default:
		InternalServerErrorResponse(w, r, err)
	}
}
