package services

import (
	"net/http"
	"web_blog/cmd/main/utils"
	"web_blog/internal/data/entity"
	"web_blog/internal/data/storage"
)

type UserService struct {
	Storage *storage.Storage
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
func (service *UserService) FindAllUsers(w http.ResponseWriter, r *http.Request) {
	var filter storage.FilterQuery
	var users []*entity.User
	var err error

	filter = storage.FilterQuery{
		Limit:  20,
		Offset: 0,
	}

	if err = filter.Parse(r); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	if err = utils.ValidateStruct(filter); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	if users, err = service.Storage.Users.FindAll(r.Context(), nil, filter); err != nil {
		utils.SwitchInternalServerErrorResponse(w, r, err)
	}

	utils.WriteJsonData(w, http.StatusOK, users)
}
