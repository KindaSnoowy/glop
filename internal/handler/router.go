package router

import (
	"blog_api/internal/handler/api"
	"blog_api/internal/handler/web"
	"blog_api/internal/repository"
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(db *sql.DB) (http.Handler, error) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

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
	web.StartWebRoutes(r, webPostHandler)

	fileServer := http.FileServer(http.Dir("../../static"))
	r.Mount("/static", http.StripPrefix("/static", fileServer))

	return r, nil

}
