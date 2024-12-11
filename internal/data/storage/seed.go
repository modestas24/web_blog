package storage

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"web_blog/internal/data/entity"

	"golang.org/x/crypto/bcrypt"
)

const (
	userURL    string = "https://randomuser.me/api/?page=3&results=200&inc=email,login"
	postURL    string = "https://dummyjson.com/posts?limit=0"
	commentURL string = "https://dummyjson.com/comments?limit=0"

	userAmount             int = 200
	postAmount             int = 400
	commentAmmount         int = 600
	userVerificationAmount int = userAmount / 2
)

type UserResponsePayload struct {
	Results []struct {
		Email string `json:"email"`
		Login struct {
			Username string `json:"username"`
			Password string `json:"password"`
		} `json:"login"`
	} `json:"results"`
}

type PostResponsePayload struct {
	Posts []struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	} `json:"posts"`
	Total int `json:"total"`
}

type CommentResponsePayload struct {
	Comments []struct {
		Body string `json:"body"`
	} `json:"comments"`
	Total int `json:"total"`
}

func processPayload(url string, payload any, handler func(*http.Response) error) error {
	var response *http.Response
	var err error

	if response, err = http.Get(url); err != nil {
		return err
	}

	if err = json.NewDecoder(response.Body).Decode(payload); err != nil {
		return err
	}

	return handler(response)
}

func SeedUsers(repository IUserRepository) error {
	var userPayload *UserResponsePayload
	return processPayload(userURL, &userPayload, func(response *http.Response) error {
		for _, u := range userPayload.Results {
			hash, err := bcrypt.GenerateFromPassword([]byte(u.Login.Password), bcrypt.DefaultCost)
			if err != nil {
				return err
			}

			user := &entity.User{
				RoleID:   1,
				Email:    u.Email,
				Username: u.Login.Username,
				Password: entity.Password{
					Hash: hash,
					Raw:  u.Login.Password,
				},
			}

			if err := repository.Create(context.TODO(), nil, user); err != nil {
				return err
			}
		}

		return nil
	})
}

func SeedPosts(repository IPostRepository) error {
	var postPayload *PostResponsePayload
	return processPayload(postURL, &postPayload, func(response *http.Response) error {
		for i := 0; i < postAmount; i++ {
			p := postPayload.Posts[i%len(postPayload.Posts)]

			post := &entity.Post{
				UserID:  rand.Int63n(int64(userAmount-2)) + 1,
				Title:   p.Title,
				Content: p.Body,
			}

			if err := repository.Create(response.Request.Context(), nil, post); err != nil {
				return err
			}
		}

		return nil
	})
}

func SeedComments(repository ICommentRepository) error {
	var commentPayload *CommentResponsePayload
	return processPayload(commentURL, &commentPayload, func(response *http.Response) error {
		for i := 0; i < commentAmmount; i++ {
			c := commentPayload.Comments[i%len(commentPayload.Comments)]

			comment := &entity.Comment{
				UserID:  rand.Int63n(int64(userAmount-2)) + 1,
				PostID:  rand.Int63n(int64(postAmount-2)) + 1,
				Content: c.Body,
			}

			if err := repository.Create(response.Request.Context(), nil, comment); err != nil {
				return err
			}
		}

		return nil
	})
}

func SeedUserVerifications(repository IVerificationRepository) error {
	for i := 0; i < commentAmmount; i++ {
		if err := repository.Create(context.TODO(), nil, &entity.Verification{
			UserID: rand.Int63n(int64(userAmount-2)) + 1,
		}); err != nil {
			return err
		}
	}
	return nil
}

func Seed(storage *Storage, count ...int) error {
	var err error
	if err = SeedUsers(storage.Users); err != nil {
		return err
	}

	if err = SeedPosts(storage.Posts); err != nil {
		return err
	}

	if err = SeedComments(storage.Comments); err != nil {
		return err
	}

	// if err = SeedUserVerifications(storage.UserVerifications); err != nil {
	// 	return err
	// }

	return nil
}
