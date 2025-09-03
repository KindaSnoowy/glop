package web

import (
	"github.com/go-chi/chi/v5"
)

func StartWebRoutes(r *chi.Mux, postHandler *PostHandler, pageHandler *PageHandler, loginHandler *LoginHandler) *chi.Mux {
	r.Route("/", func(r chi.Router) {
		r.Get("/", pageHandler.GetHomePage)
	})

	r.Route("/login", func(r chi.Router) {
		r.Get("/", loginHandler.GetLoginPage)
		r.Post("/", loginHandler.Login)
	})

	r.Route("/posts", func(r chi.Router) {
		r.Get("/", postHandler.GetPostsPage)
		r.Get("/{id}", postHandler.GetPostIdPage)
	})

	return r
}
