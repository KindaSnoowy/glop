package web

import (
	"github.com/go-chi/chi/v5"
)

func StartWebRoutes(r *chi.Mux, postHandler *PostHandler) *chi.Mux {
	r.Route("/posts", func(r chi.Router) {
		r.Get("/", postHandler.GetPostsPage)
		r.Get("/{id}", postHandler.GetPostIdPage)
	})

	return r
}
