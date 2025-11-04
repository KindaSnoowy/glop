// Package web -> todas as rotas web, que servem diretamente ao site acessado
package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func StartWebRoutes(r *chi.Mux, authMiddleware func(http.Handler) http.Handler, permissionMiddleware func(http.Handler) http.Handler,
	postHandler *PostHandler, homeHandler *HomeHandler, loginHandler *LoginHandler, newPostHandler *NewPostHandler,
) *chi.Mux {
	r.Route("/", func(r chi.Router) {
		r.Get("/", homeHandler.GetHomePage)
	})

	r.Route("/posts", func(r chi.Router) {
		r.Get("/", postHandler.GetPostsPage)
		r.Get("/{id}", postHandler.GetPostIDPage)

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)
			r.Use(permissionMiddleware)

			r.Get("/edit", newPostHandler.GetNewPostPage)
			r.Get("/edit/{id}", newPostHandler.GetNewPostPage)

			// rotas para criar e editar os posts (api do web)
			r.Post("/", postHandler.WebCreatePost)
			r.Put("/{id}", postHandler.WebUpdatePost)
		})
	})

	r.Route("/login", func(r chi.Router) {
		r.Get("/", loginHandler.GetLoginPage)
		r.Post("/", loginHandler.WebLogin)
	})

	r.Get("/logout", loginHandler.LogoutHandler)
	return r
}
