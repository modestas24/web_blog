package services

import (
	"net/http"
)

type IHealthService interface {
	CheckHealth(w http.ResponseWriter, _ *http.Request)
}

type IAuthenticationService interface {
	RegisterUser(http.ResponseWriter, *http.Request)
	VerifyUser(http.ResponseWriter, *http.Request)
	LoginUser(http.ResponseWriter, *http.Request)
	LogoutUser(http.ResponseWriter, *http.Request)
}

type IUserService interface {
	FindAllUsers(http.ResponseWriter, *http.Request)
}

type IPostService interface {
	CreatePost(http.ResponseWriter, *http.Request)
	FindAllPosts(http.ResponseWriter, *http.Request)
	FindAllPostsByUserID(http.ResponseWriter, *http.Request)
	FindPost(http.ResponseWriter, *http.Request)
	UpdatePost(http.ResponseWriter, *http.Request)
	DeletePost(http.ResponseWriter, *http.Request)
}

type ICommentService interface {
	CreateComment(http.ResponseWriter, *http.Request)
	FindAllComments(http.ResponseWriter, *http.Request)
	FindAllCommentsByPostID(http.ResponseWriter, *http.Request)
	DeleteComment(http.ResponseWriter, *http.Request)
}

type Services struct {
	Health  IHealthService
	Auth    IAuthenticationService
	User    IUserService
	Post    IPostService
	Comment ICommentService
}
