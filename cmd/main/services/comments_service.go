package services

import (
	"net/http"
	"strconv"
	"web_blog/cmd/main/middlewares"
	"web_blog/cmd/main/utils"
	"web_blog/internal/data/entity"
	"web_blog/internal/data/storage"

	"github.com/go-chi/chi/v5"
)

type CommentService struct {
	Storage *storage.Storage
}

type CreateCommentPayload struct {
	Content string `json:"content" service.Application.Validator:"required,max=512"`
}

// CreateComment godoc
//
//	@Summary		Create a comment
//	@Description	Add a comment to a specific post
//	@Tags			comments
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"Post ID"
//	@Param			payload	body		CreateCommentPayload	true	"Comment payload"
//	@Success		201		{object}	EnvelopeJson{data=entity.Comment}
//	@Failure		400		{object}	ErrorEnvelopeJson
//	@Failure		500		{object}	ErrorEnvelopeJson
//	@Security		ApiKeyAuth
//	@Router			/posts/{id}/comments [post]
func (service *CommentService) CreateComment(w http.ResponseWriter, r *http.Request) {
	var comment *entity.Comment
	var payload CreateCommentPayload
	var id int
	var err error

	if id, err = strconv.Atoi(chi.URLParam(r, "id")); err != nil {
		utils.InternalServerErrorResponse(w, r, err)
		return
	}

	if err = utils.ReadJson(w, r, &payload); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	if err = utils.ValidateStruct(payload); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	comment = &entity.Comment{
		PostID:  int64(id),
		UserID:  middlewares.FindUserFromContext(r).ID,
		Content: payload.Content,
	}

	if err = service.Storage.Comments.Create(r.Context(), nil, comment); err != nil {
		utils.SwitchInternalServerErrorResponse(w, r, err)
		return
	}

	utils.WriteJsonData(w, http.StatusCreated, comment)
}

// FindAllComments godoc
//
//	@Summary		Get all comments
//	@Description	Retrieve all comments in the system
//	@Tags			comments
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int	false	"Limit"
//	@Param			offset	query		int	false	"Offset"
//	@Success		200		{array}		EnvelopeJson{data=[]entity.Comment}
//	@Failure		500		{object}	ErrorEnvelopeJson
//	@Security		ApiKeyAuth
//	@Router			/posts/comments [get]
func (service *CommentService) FindAllComments(w http.ResponseWriter, r *http.Request) {
	var filter storage.FilterQuery
	var comments []*entity.Comment
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

	if comments, err = service.Storage.Comments.FindAll(r.Context(), nil, filter); err != nil {
		utils.SwitchInternalServerErrorResponse(w, r, err)
		return
	}

	utils.WriteJsonData(w, http.StatusOK, comments)
}

// FindAllCommentsByPostID godoc
//
//	@Summary		Get all comments by post ID
//	@Description	Retrieve all comments associated with a specific post
//	@Tags			comments
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int	true	"Post ID"
//	@Param			limit	query		int	false	"Limit"
//	@Param			offset	query		int	false	"Offset"
//	@Success		200		{array}		EnvelopeJson{data=[]entity.Comment}
//	@Failure		400		{object}	ErrorEnvelopeJson
//	@Failure		500		{object}	ErrorEnvelopeJson
//	@Router			/posts/{id}/comments [get]
func (service *CommentService) FindAllCommentsByPostID(w http.ResponseWriter, r *http.Request) {
	var filter storage.FilterQuery
	var comments []*entity.Comment
	var id int
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

	if id, err = strconv.Atoi(chi.URLParam(r, "id")); err != nil {
		utils.InternalServerErrorResponse(w, r, err)
		return
	}

	if comments, err = service.Storage.Comments.FindAllByPostID(r.Context(), nil, filter, int64(id)); err != nil {
		utils.SwitchInternalServerErrorResponse(w, r, err)
		return
	}

	utils.WriteJsonData(w, http.StatusOK, comments)
}

// DeleteComment godoc
//
//	@Summary		Delete a comment
//	@Description	Delete a specific comment by its ID
//	@Tags			comments
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Comment ID"
//	@Success		204	"No Content"
//	@Failure		400	{object}	ErrorEnvelopeJson
//	@Failure		404	{object}	ErrorEnvelopeJson
//	@Failure		500	{object}	ErrorEnvelopeJson
//	@Security		ApiKeyAuth
//	@Router			/posts/comments/{id} [delete]
func (service *CommentService) DeleteComment(w http.ResponseWriter, r *http.Request) {
	var id int
	var err error

	if id, err = strconv.Atoi(chi.URLParam(r, "id")); err != nil {
		utils.InternalServerErrorResponse(w, r, err)
		return
	}

	if err = service.Storage.Comments.Delete(r.Context(), nil, int64(id)); err != nil {
		utils.SwitchInternalServerErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
