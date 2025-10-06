package web

import (
	"html/template"
	"log"
	"net/http"
)

type LoginHandler struct{}

func StartLoginHandler() *LoginHandler {
	return &LoginHandler{}
}

func (s *LoginHandler) GetLoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../../view/pages/login.html")
	if err != nil {
		log.Printf("Erro ao carregar página de login: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Erro ao executar template: %v", err)
	}
}
