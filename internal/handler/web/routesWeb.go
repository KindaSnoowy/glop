// Package web -> todas as rotas web, que servem diretamente ao site acessado
package web

import (
	"github.com/go-chi/chi/v5"
)

func StartWebRoutes(r *chi.Mux, postHandler *PostHandler, homeHandler *HomeHandler,
	loginHandler *LoginHandler,
) *chi.Mux {
	r.Route("/", func(r chi.Router) {
		r.Get("/", homeHandler.GetHomePage)
	})

	r.Route("/posts", func(r chi.Router) {
		r.Get("/", postHandler.GetPostsPage)
		r.Get("/{id}", postHandler.GetPostIDPage)
	})

	r.Route("/login", func(r chi.Router) {
		r.Get("/", loginHandler.GetLoginPage)
	})
	return r
}
