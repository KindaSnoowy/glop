// Package api -> todas as rotas da API, não retornam html
package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func StartAPIRoutes(r *chi.Mux,
	authMiddleware func(http.Handler) http.Handler, permissionMiddleware func(http.Handler) http.Handler,
	postHandler *PostHandlerAPI, userHandler *UserHandlerAPI, loginHandler *LoginHandlerAPI,
) *chi.Mux {
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
			r.Get("/{id}", userHandler.GetUserByID)
			r.Get("/username/{username}", userHandler.GetUserByUsername)

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware)
				r.Use(permissionMiddleware)
				r.Put("/{id}", userHandler.UpdateUser)
				r.Post("/", userHandler.CreateUser)

			})
		})
	})

	return r
}
