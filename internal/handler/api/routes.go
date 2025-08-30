package api

import (
	"github.com/go-chi/chi/v5"
)

func StartApiRoutes(r *chi.Mux, postHandler *PostHandler_api) *chi.Mux {
	// posts
	r.Route("/api", func(r chi.Router) {
		r.Route("/posts", func(r chi.Router) {
			r.Get("/", postHandler.GetPosts)
			r.Post("/", postHandler.CreatePost)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", postHandler.GetPostById)
				r.Put("/", postHandler.Update)
				r.Delete("/", postHandler.DeletePost)
			})
		})

	})

	return r
}
