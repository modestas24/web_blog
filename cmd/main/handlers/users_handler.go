package handlers

import (
	"net/http"
	"web_blog/internal/data/entity"
	"web_blog/internal/data/storage"
)

type UserHandler struct {
	Storage      *storage.Storage
	ErrorHandler IErrorHandler
}

type userKey string

const UserCtx userKey = "user"

func findUserFromCtx(r *http.Request) *entity.User {
	user, _ := r.Context().Value(UserCtx).(*entity.User)
	return user
}

// FindAllUsers godoc
//
//	@Summary		Get all users
//	@Description	Retrieve a list of all registered users
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int	false	"Limit"
//	@Param			offset	query		int	false	"Offset"
//	@Success		200		{object}	EnvelopeJson{data=[]entity.User}
//	@Failure		500		{object}	ErrorEnvelopeJson
//	@Security		ApiKeyAuth
//	@Router			/users [get]
func (handler *UserHandler) FindAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	var filter storage.FilterQuery
	var users []*entity.User
	var err error

	filter = storage.FilterQuery{
		Limit:  20,
		Offset: 0,
	}

	if err = filter.Parse(r); err != nil {
		handler.ErrorHandler.BadRequestError(w, r, err)
		return
	}

	if err = validateStruct(filter); err != nil {
		handler.ErrorHandler.BadRequestError(w, r, err)
		return
	}

	if users, err = handler.Storage.Users.FindAll(r.Context(), nil, filter); err != nil {
		handler.ErrorHandler.SwitchInternalServerError(w, r, err)
	}

	writeJsonData(w, http.StatusOK, users)
}
