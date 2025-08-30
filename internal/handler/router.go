package router

import (
	"blog_api/internal/handler/api"
	"blog_api/internal/handler/web"
	"blog_api/internal/repository"
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(db *sql.DB) (http.Handler, error) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println(w, r)
			h.ServeHTTP(w, r)
		})
	})

	// inicializa rotas API
	// inicializa repository posts
	postRepo, err := repository.StartPostRepository(db)
	if err != nil {
		return nil, err
	}
	apiPostHandler := api.StartPostHandler(postRepo)
	api.StartApiRoutes(r, apiPostHandler)

	// inicializa rotas WEB
	webPostHandler := web.StartPostHandler(postRepo)
	webHomeHandler := web.StartHomeHandler()
	web.StartWebRoutes(r, webPostHandler, webHomeHandler)

	// acesso ao static
	fileServer := http.FileServer(http.Dir("../../static"))
	r.Mount("/static", http.StripPrefix("/static", fileServer))

	return r, nil

}
