package handlers

import (
	"context"
	"net/http"
	"strconv"
	"web_blog/internal/data/entity"
	"web_blog/internal/data/storage"

	"github.com/go-chi/chi/v5"
)

type PostHandler struct {
	Storage      *storage.Storage
	ErrorHandler IErrorHandler
}

type postKey string

const PostCtx postKey = "post"

func (handler *PostHandler) FindPostFromContext(r *http.Request) *entity.Post {
	post, _ := r.Context().Value(PostCtx).(*entity.Post)
	return post
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
func (handler *PostHandler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	var err error

	if err = readJson(w, r, &payload); err != nil {
		handler.ErrorHandler.InternalServerError(w, r, err)
		return
	}

	if err = validateStruct(payload); err != nil {
		handler.ErrorHandler.BadRequestError(w, r, err)
		return
	}

	post := &entity.Post{
		Title:   payload.Title,
		Content: payload.Content,
		UserID:  findUserFromCtx(r).ID,
	}

	if err = handler.Storage.Posts.Create(r.Context(), nil, post); err != nil {
		handler.ErrorHandler.SwitchInternalServerError(w, r, err)
		return
	}

	writeJsonData(w, http.StatusCreated, post)
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
func (handler *PostHandler) FindAllPostHandler(w http.ResponseWriter, r *http.Request) {
	var filter storage.FilterQuery
	var posts []*entity.Post
	var err error
	ctx := r.Context()

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

	if posts, err = handler.Storage.Posts.FindAll(ctx, nil, filter); err != nil {
		handler.ErrorHandler.SwitchInternalServerError(w, r, err)
		return
	}

	writeJsonData(w, http.StatusOK, posts)
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
func (handler *PostHandler) FindAllPostByUserIDHandler(w http.ResponseWriter, r *http.Request) {
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
		handler.ErrorHandler.BadRequestError(w, r, err)
		return
	}

	if err = validateStruct(filter); err != nil {
		handler.ErrorHandler.BadRequestError(w, r, err)
		return
	}

	if id, err = strconv.Atoi(chi.URLParam(r, "id")); err != nil {
		handler.ErrorHandler.InternalServerError(w, r, err)
		return
	}

	if posts, err = handler.Storage.Posts.FindAllByUserID(ctx, nil, filter, int64(id)); err != nil {
		handler.ErrorHandler.SwitchInternalServerError(w, r, err)
		return
	}

	writeJsonData(w, http.StatusOK, posts)
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
func (handler *PostHandler) FindPostHandler(w http.ResponseWriter, r *http.Request) {
	post := handler.FindPostFromContext(r)
	writeJsonData(w, http.StatusOK, post)
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
func (handler *PostHandler) UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := handler.FindPostFromContext(r)
	var payload UpdatePostPayload
	var err error

	if err = readJson(w, r, &payload); err != nil {
		handler.ErrorHandler.BadRequestError(w, r, err)
		return
	}

	if err = validateStruct(payload); err != nil {
		handler.ErrorHandler.BadRequestError(w, r, err)
		return
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}

	if err = handler.Storage.Posts.Update(r.Context(), nil, post); err != nil {
		handler.ErrorHandler.SwitchInternalServerError(w, r, err)
		return
	}

	writeJsonData(w, http.StatusOK, post)
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
func (handler *PostHandler) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var id int
	var err error

	if id, err = strconv.Atoi(chi.URLParam(r, "id")); err != nil {
		handler.ErrorHandler.InternalServerError(w, r, err)
		return
	}

	if err = handler.Storage.Posts.Delete(ctx, nil, int64(id)); err != nil {
		handler.ErrorHandler.SwitchInternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (handler *PostHandler) PostContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		var id int
		var post *entity.Post
		var err error

		if id, err = strconv.Atoi(chi.URLParam(r, "id")); err != nil {
			handler.ErrorHandler.InternalServerError(w, r, err)
			return
		}

		if post, err = handler.Storage.Posts.Find(ctx, nil, int64(id)); err != nil {
			handler.ErrorHandler.SwitchInternalServerError(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, PostCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
