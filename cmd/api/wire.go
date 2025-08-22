package main

import (
	"blog_api/internal/handlers"
	"blog_api/internal/repository"
	"database/sql"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

func bootstrap(db *sql.DB) (http.Handler, error) {
	// inicializa tabela posts
	createTableSQL := `CREATE TABLE IF NOT EXISTS posts (
			"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, "title" TEXT, "content" TEXT,
			"createdAt" DATETIME, "updatedAt" DATETIME
		);`
	if _, err := db.Exec(createTableSQL); err != nil {
		return nil, err
	}
	log.Println("Tabela 'posts' pronta")

	// inicializando posts handler
	postRepo := repository.StartPostRepository(db)
	postHandler := handlers.StartPostHandler(postRepo)

	// inicializa o router
	router := handlers.NewRouter(postHandler)

	return router, nil
}
