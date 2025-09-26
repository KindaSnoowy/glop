package web

import (
	"log"
	"net/http"
	"text/template"
)

type HomeHandler struct {
}

func StartHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (s *HomeHandler) GetHomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../../view/pages/home.html")
	if err != nil {
		log.Printf("Erro ao carregar template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Erro ao executar template: %v", err)
	}
}
