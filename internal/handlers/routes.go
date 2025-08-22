package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(postHandler *PostHandler) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/posts", func(r chi.Router) {
		r.Get("/", postHandler.GetPosts)
		r.Post("/", postHandler.CreatePost)
	})

	return r
}
