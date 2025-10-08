// Package router -> Injeta todas as dependencias no router
package router

import (
	"database/sql"
	"log"
	"net/http"

	"blog_api/internal/handler/api"
	"blog_api/internal/handler/web"
	middlewares "blog_api/internal/middleware"
	"blog_api/internal/repository"

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

	// testando commit
	// inicializa rotas API
	// inicializa repositories
	postRepo, err := repository.StartPostRepository(db)
	if err != nil {
		return nil, err
	}
	userRepo, err := repository.StartUserRepository(db)
	if err != nil {
		return nil, err
	}
	sessionRepo, err := repository.StartSessionRepository(db)
	if err != nil {
		return nil, err
	}

	// inicializa handlers
	apiPostHandler := api.StartPostHandler(postRepo)
	apiLoginHandler := api.StartLoginHandler(sessionRepo, userRepo)
	apiUserHandler := api.StartUserHandler(userRepo)

	authMiddleware := middlewares.AuthMiddleware(sessionRepo)
	permissionMiddleware := middlewares.PermissionMiddleware(userRepo)
	api.StartAPIRoutes(r, authMiddleware, permissionMiddleware, apiPostHandler, apiUserHandler, apiLoginHandler)

	// inicializa rotas WEB
	webPostHandler := web.StartPostHandler(postRepo)
	webLoginHandler := web.StartLoginHandler()
	webHomeHandler := web.StartHomeHandler()
	web.StartWebRoutes(r, webPostHandler, webHomeHandler, webLoginHandler)

	// acesso ao static
	fileServer := http.FileServer(http.Dir("../../static"))
	r.Mount("/static", http.StripPrefix("/static", fileServer))

	return r, nil
}
