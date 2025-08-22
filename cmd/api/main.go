package main

import (
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

	router, err := bootstrap(db)
	if err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
