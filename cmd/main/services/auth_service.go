package services

import (
	"net/http"
	"web_blog/cmd/main/utils"

	"web_blog/internal/authentication"
	"web_blog/internal/data/entity"
	"web_blog/internal/data/storage"

	"github.com/google/uuid"
)

type AuthService struct {
	Storage       *storage.Storage
	Authenticator *authentication.StatefulAuthenticator
}

type RegisterUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

// RegisterUser godoc
//
//	@Summary		Register a new user
//	@Description	Create a new user account with the provided details
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"Registration details"
//	@Success		201		{object}	EnvelopeJson{data=entity.User}
//	@Failure		400		{object}	ErrorEnvelopeJson
//	@Failure		500		{object}	ErrorEnvelopeJson
//	@Router			/authentication/register [post]
func (service *AuthService) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user *entity.User
	var password entity.Password
	var payload RegisterUserPayload
	var err error

	if err = utils.ReadJson(w, r, &payload); err != nil {
		utils.InternalServerErrorResponse(w, r, err)
		return
	}

	if err = utils.ValidateStruct(payload); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	password = entity.Password{}
	if err = password.Set(payload.Password); err != nil {
		utils.InternalServerErrorResponse(w, r, err)
		return
	}

	user = &entity.User{
		RoleID:   1,
		Email:    payload.Email,
		Username: payload.Username,
		Password: password,
	}

	if err = service.Storage.Users.CreateWithVerification(
		r.Context(),
		nil,
		service.Storage.Verifications,
		user,
	); err != nil {
		utils.SwitchInternalServerErrorResponse(w, r, err)
		return
	}

	utils.WriteJsonData(w, http.StatusCreated, user)
}

type VerifyUserPayload struct {
	UUID uuid.UUID `json:"id" validate:"uuid"`
}

// VerifyUser godoc
//
//	@Summary		Verify a user account
//	@Description	Verify a user using a UUID from an email or verification method
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		VerifyUserPayload	true	"Verification details"
//	@Success		200		{object}	EnvelopeJson{data=entity.User}
//	@Failure		400		{object}	ErrorEnvelopeJson
//	@Failure		500		{object}	ErrorEnvelopeJson
//	@Router			/authentication/verify [post]
func (service *AuthService) VerifyUser(w http.ResponseWriter, r *http.Request) {
	var user *entity.User
	var payload VerifyUserPayload
	var err error

	if err = utils.ReadJson(w, r, &payload); err != nil {
		utils.InternalServerErrorResponse(w, r, err)
		return
	}

	if err = utils.ValidateStruct(payload); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	user = &entity.User{}

	if err = service.Storage.Users.Verify(
		r.Context(),
		nil,
		service.Storage.Verifications,
		payload.UUID,
		user,
	); err != nil {
		utils.SwitchInternalServerErrorResponse(w, r, err)
		return
	}

	utils.WriteJsonData(w, http.StatusOK, user)
}

type TokenEnvelopeJson struct {
	Token string `json:"token"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

// LoginUser godoc
//
//	@Summary		User login
//	@Description	Authenticate a user and return a session token
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		LoginUserPayload	true	"Login details"
//	@Success		202		{object}	EnvelopeJson{data=TokenEnvelopeJson}
//	@Failure		400		{object}	ErrorEnvelopeJson
//	@Failure		500		{object}	ErrorEnvelopeJson
//	@Router			/authentication/login [post]
func (service *AuthService) LoginUser(w http.ResponseWriter, r *http.Request) {
	var token string
	var user *entity.User
	var payload *LoginUserPayload
	var err error

	payload = &LoginUserPayload{}

	if err = utils.ReadJson(w, r, &payload); err != nil {
		utils.InternalServerErrorResponse(w, r, err)
		return
	}

	if err = utils.ValidateStruct(payload); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	if user, err = service.Storage.Users.FindByEmail(r.Context(), nil, payload.Email); err != nil {
		utils.InternalServerErrorResponse(w, r, err)
		return
	}

	if err = user.Password.Compare([]byte(payload.Password)); err != nil {
		utils.InternalServerErrorResponse(w, r, err)
		return
	}

	if token, err = service.Authenticator.Generate(); err != nil {
		utils.InternalServerErrorResponse(w, r, err)
		return
	}

	if err = service.Authenticator.Create(r.Context(), service.Storage.Sessions, token, user.ID); err != nil {
		utils.InternalServerErrorResponse(w, r, err)
		return
	}

	utils.WriteJsonData(w, http.StatusAccepted, TokenEnvelopeJson{Token: token})
}

type LogoutUserPayload struct {
	Token string `json:"token" validate:"required"`
}

// LogoutUser godoc
//
//	@Summary		User logout
//	@Description	Invalidate a session token for a user
//	@Tags			authentication
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body	LogoutUserPayload	true	"Logout details"
//	@Success		204		"No Content"
//	@Failure		400		{object}	ErrorEnvelopeJson
//	@Failure		404		{object}	ErrorEnvelopeJson
//	@Failure		500		{object}	ErrorEnvelopeJson
//	@Router			/authentication/logout [delete]
func (service *AuthService) LogoutUser(w http.ResponseWriter, r *http.Request) {
	var payload LogoutUserPayload
	var err error

	if err = utils.ReadJson(w, r, &payload); err != nil {
		utils.InternalServerErrorResponse(w, r, err)
		return
	}

	if err = utils.ValidateStruct(payload); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	if err = service.Authenticator.Invalidate(r.Context(), service.Storage.Sessions, payload.Token); err != nil {
		utils.NotFoundResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
