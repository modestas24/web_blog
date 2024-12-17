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

type PostService struct {
	Storage *storage.Storage
}

type CreatePostPayload struct {
	Title   string `json:"title" validate:"required,max=128"`
	Content string `json:"content" validate:"required,max=1024"`
}

// CreatePost godoc
//
//	@Summary		Create a new post
//	@Description	Create a new post with the given payload
//	@Tags			posts
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreatePostPayload	true	"Post payload"
//	@Success		201		{object}	EnvelopeJson{data=entity.Post}
//	@Failure		400		{object}	ErrorEnvelopeJson
//	@Failure		500		{object}	ErrorEnvelopeJson
//	@Router			/posts [post]
func (service *PostService) CreatePost(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	var err error

	if err = utils.ReadJson(w, r, &payload); err != nil {
		utils.InternalServerErrorResponse(w, r, err)
		return
	}

	if err = utils.ValidateStruct(payload); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	post := &entity.Post{
		Title:   payload.Title,
		Content: payload.Content,
		UserID:  middlewares.FindUserFromContext(r).ID,
	}

	if err = service.Storage.Posts.Create(r.Context(), nil, post); err != nil {
		utils.SwitchInternalServerErrorResponse(w, r, err)
		return
	}

	utils.WriteJsonData(w, http.StatusCreated, post)
}

// FindAllPosts godoc
//
//	@Summary		Get all posts
//	@Description	Retrieve a list of all posts
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int	false	"Limit"
//	@Param			offset	query		int	false	"Offset"
//	@Success		200		{array}		EnvelopeJson{data=entity.Post}
//	@Failure		500		{object}	ErrorEnvelopeJson
//	@Router			/posts [get]
func (service *PostService) FindAllPosts(w http.ResponseWriter, r *http.Request) {
	var filter storage.FilterQuery
	var posts []*entity.Post
	var err error
	ctx := r.Context()

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

	if posts, err = service.Storage.Posts.FindAll(ctx, nil, filter); err != nil {
		utils.SwitchInternalServerErrorResponse(w, r, err)
		return
	}

	utils.WriteJsonData(w, http.StatusOK, posts)
}

// FindAllPostsByUserID godoc
//
//	@Summary		Get posts by user ID
//	@Description	Retrieve all posts created by a specific user
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int	true	"User ID"
//	@Param			limit	query		int	false	"Limit"
//	@Param			offset	query		int	false	"Offset"
//	@Success		200		{array}		EnvelopeJson{data=entity.Post}
//	@Failure		400		{object}	ErrorEnvelopeJson
//	@Failure		500		{object}	ErrorEnvelopeJson
//	@Router			/posts/user/{id} [get]
func (service *PostService) FindAllPostsByUserID(w http.ResponseWriter, r *http.Request) {
	var filter storage.FilterQuery
	var posts []*entity.Post
	var id int
	var err error
	ctx := r.Context()

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

	if posts, err = service.Storage.Posts.FindAllByUserID(ctx, nil, filter, int64(id)); err != nil {
		utils.SwitchInternalServerErrorResponse(w, r, err)
		return
	}

	utils.WriteJsonData(w, http.StatusOK, posts)
}

// FindPost godoc
//
//	@Summary		Get a post by ID
//	@Description	Retrieve a specific post by its ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Post ID"
//	@Success		200	{object}	EnvelopeJson{data=entity.Post}
//	@Failure		404	{object}	ErrorEnvelopeJson
//	@Failure		500	{object}	ErrorEnvelopeJson
//	@Router			/posts/{id} [get]
func (service *PostService) FindPost(w http.ResponseWriter, r *http.Request) {
	post := middlewares.FindPostFromContext(r)
	utils.WriteJsonData(w, http.StatusOK, post)
}

type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=128"`
	Content *string `json:"content" validate:"omitempty,max=1024"`
}

// UpdatePost godoc
//
//	@Summary		Update a post
//	@Description	Update the details of a specific post
//	@Tags			posts
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Post ID"
//	@Param			payload	body		UpdatePostPayload	true	"Update payload"
//	@Success		200		{object}	EnvelopeJson{data=entity.Post}
//	@Failure		400		{object}	ErrorEnvelopeJson
//	@Failure		404		{object}	ErrorEnvelopeJson
//	@Failure		500		{object}	ErrorEnvelopeJson
//	@Router			/posts/{id} [patch]
func (service *PostService) UpdatePost(w http.ResponseWriter, r *http.Request) {
	post := middlewares.FindPostFromContext(r)
	var payload UpdatePostPayload
	var err error

	if err = utils.ReadJson(w, r, &payload); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	if err = utils.ValidateStruct(payload); err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}

	if err = service.Storage.Posts.Update(r.Context(), nil, post); err != nil {
		utils.SwitchInternalServerErrorResponse(w, r, err)
		return
	}

	utils.WriteJsonData(w, http.StatusOK, post)
}

// DeletePost godoc
//
//	@Summary		Delete a post
//	@Description	Delete a specific post by its ID
//	@Tags			posts
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Post ID"
//	@Success		204	"No Content"
//	@Failure		400	{object}	ErrorEnvelopeJson
//	@Failure		404	{object}	ErrorEnvelopeJson
//	@Failure		500	{object}	ErrorEnvelopeJson
//	@Router			/posts/{id} [delete]
func (service *PostService) DeletePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var id int
	var err error

	if id, err = strconv.Atoi(chi.URLParam(r, "id")); err != nil {
		utils.InternalServerErrorResponse(w, r, err)
		return
	}

	if err = service.Storage.Posts.Delete(ctx, nil, int64(id)); err != nil {
		utils.SwitchInternalServerErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
