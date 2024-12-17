package middlewares

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"web_blog/cmd/main/utils"
	"web_blog/internal/data/entity"
)

type postKey string

const PostCtx postKey = "post"

func (middleware *Middleware) PostContext(next http.Handler) http.Handler {
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
			utils.InternalServerErrorResponse(w, r, err)
			return
		}

		if post, err = middleware.Storage.Posts.Find(ctx, nil, int64(id)); err != nil {
			utils.SwitchInternalServerErrorResponse(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, PostCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func FindPostFromContext(r *http.Request) *entity.Post {
	post, _ := r.Context().Value(PostCtx).(*entity.Post)
	return post
}
