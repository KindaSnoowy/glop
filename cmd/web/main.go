package main

import (
	router "blog_api/internal/handler"
	"database/sql"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "../../glop.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("Database connection success")

	router, err := router.NewRouter(db)
	if err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(":7777", router); err != nil {
		log.Fatal(err)
	}
}
