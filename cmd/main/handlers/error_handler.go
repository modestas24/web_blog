package handlers

import (
	"errors"
	"net/http"
	"web_blog/internal/data/storage"

	"github.com/jackc/pgx"
	"go.uber.org/zap"
)

type ErrorHandler struct {
	Logger *zap.Logger
}

func (handler *ErrorHandler) InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	handler.Logger.Warn("internal server error", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.Error(err))
	writeJsonError(w, r, http.StatusInternalServerError, err.Error())
}

func (handler *ErrorHandler) BadRequestError(w http.ResponseWriter, r *http.Request, err error) {
	handler.Logger.Warn("bad request error", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.Error(err))
	writeJsonError(w, r, http.StatusBadRequest, err.Error())
}

func (handler *ErrorHandler) NotFoundError(w http.ResponseWriter, r *http.Request, err error) {
	handler.Logger.Warn("not found error", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.Error(err))
	writeJsonError(w, r, http.StatusNotFound, err.Error())
}

func (handler *ErrorHandler) ForbiddenError(w http.ResponseWriter, r *http.Request, err error) {
	handler.Logger.Warn("forbidden error", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.Error(err))
	writeJsonError(w, r, http.StatusForbidden, err.Error())
}

func (handler *ErrorHandler) UnauthorizedError(w http.ResponseWriter, r *http.Request, err error) {
	handler.Logger.Warn("unauthorized error", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.Error(err))
	writeJsonError(w, r, http.StatusUnauthorized, err.Error())
}

func (handler *ErrorHandler) SwitchInternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, storage.ErrorNotFound):
		handler.NotFoundError(w, r, err)
		return
	case errors.Is(err, pgx.ErrNoRows):
		handler.NotFoundError(w, r, storage.ErrorNotFound)
		return
	default:
		handler.InternalServerError(w, r, err)
	}
}
