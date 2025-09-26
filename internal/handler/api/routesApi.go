package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func StartApiRoutes(r *chi.Mux, authMiddleware func(http.Handler) http.Handler, permissionMiddleware func(http.Handler) http.Handler, postHandler *PostHandler_api, userHandler *UserHandler_api, loginHandler *LoginHandler_api) *chi.Mux {
	r.Route("/api", func(r chi.Router) {
		r.Post("/login", loginHandler.Login)

		r.Route("/posts", func(r chi.Router) {
			r.Get("/", postHandler.GetPosts)
			r.Get("/{id}", postHandler.GetPostById)

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware)
				r.Use(permissionMiddleware)
				r.Post("/", postHandler.CreatePost)
				r.Put("/{id}", postHandler.Update)
				r.Delete("/{id}", postHandler.DeletePost)
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Post("/", userHandler.CreateUser)

			r.Get("/{id}", userHandler.GetUserById)
			r.Get("/username/{username}", userHandler.GetUserByUsername)

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware)
				r.Use(permissionMiddleware)
				r.Put("/{id}", userHandler.UpdateUser)
			})
		})
	})

	return r
}
