package handlers

import (
	"net/http"

	"web_blog/internal/auth"
	"web_blog/internal/data/entity"
	"web_blog/internal/data/storage"

	"github.com/google/uuid"
)

type AuthHandler struct {
	Storage         *storage.Storage
	Authentificator *auth.StatefulAuthenticator
	ErrorHandler    IErrorHandler
}

type RegisterUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

type VerifyUserPayload struct {
	UUID uuid.UUID `json:"id" validate:"uuid"`
}

type LogoutUserPayload struct {
	Token string `json:"token" validate:"required"`
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
func (handler *AuthHandler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var user *entity.User
	var password entity.Password
	var payload RegisterUserPayload
	var err error

	if err = readJson(w, r, &payload); err != nil {
		handler.ErrorHandler.InternalServerError(w, r, err)
		return
	}

	if err = validateStruct(payload); err != nil {
		handler.ErrorHandler.BadRequestError(w, r, err)
		return
	}

	password = entity.Password{}
	if err = password.Set(payload.Password); err != nil {
		handler.ErrorHandler.InternalServerError(w, r, err)
		return
	}

	user = &entity.User{
		RoleID:   1,
		Email:    payload.Email,
		Username: payload.Username,
		Password: password,
	}

	if err = handler.Storage.Users.CreateWithVerification(
		r.Context(),
		nil,
		handler.Storage.Verifications,
		user,
	); err != nil {
		handler.ErrorHandler.SwitchInternalServerError(w, r, err)
		return
	}

	writeJsonData(w, http.StatusCreated, user)
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
func (handler *AuthHandler) VerifyUserHandler(w http.ResponseWriter, r *http.Request) {
	var user *entity.User
	var payload VerifyUserPayload
	var err error

	if err = readJson(w, r, &payload); err != nil {
		handler.ErrorHandler.InternalServerError(w, r, err)
		return
	}

	if err = validateStruct(payload); err != nil {
		handler.ErrorHandler.BadRequestError(w, r, err)
		return
	}

	user = &entity.User{}

	if err = handler.Storage.Users.Verify(
		r.Context(),
		nil,
		handler.Storage.Verifications,
		payload.UUID,
		user,
	); err != nil {
		handler.ErrorHandler.SwitchInternalServerError(w, r, err)
		return
	}

	writeJsonData(w, http.StatusOK, user)
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
func (handler *AuthHandler) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var token string
	var user *entity.User
	var payload *LoginUserPayload
	var err error

	payload = &LoginUserPayload{}

	if err = readJson(w, r, &payload); err != nil {
		handler.ErrorHandler.InternalServerError(w, r, err)
		return
	}

	if err = validateStruct(payload); err != nil {
		handler.ErrorHandler.BadRequestError(w, r, err)
		return
	}

	if user, err = handler.Storage.Users.FindByEmail(r.Context(), nil, payload.Email); err != nil {
		handler.ErrorHandler.InternalServerError(w, r, err)
		return
	}

	if err = user.Password.Compare([]byte(payload.Password)); err != nil {
		handler.ErrorHandler.InternalServerError(w, r, err)
		return
	}

	if token, err = handler.Authentificator.Generate(); err != nil {
		handler.ErrorHandler.InternalServerError(w, r, err)
		return
	}

	if err = handler.Authentificator.Create(r.Context(), handler.Storage.Sessions, token, user.ID); err != nil {
		handler.ErrorHandler.InternalServerError(w, r, err)
		return
	}

	writeJsonData(w, http.StatusAccepted, TokenEnvelopeJson{Token: token})
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
func (handler *AuthHandler) LogoutUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload LogoutUserPayload
	var err error

	if err = readJson(w, r, &payload); err != nil {
		handler.ErrorHandler.InternalServerError(w, r, err)
		return
	}

	if err = validateStruct(payload); err != nil {
		handler.ErrorHandler.BadRequestError(w, r, err)
		return
	}

	if err = handler.Authentificator.Invalidate(r.Context(), handler.Storage.Sessions, payload.Token); err != nil {
		handler.ErrorHandler.NotFoundError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
