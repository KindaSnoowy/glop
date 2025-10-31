package web

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

type HomeHandler struct{}

type homePageData struct {
	IsAuthenticated bool
}

func StartHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (s *HomeHandler) GetHomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../../view/pages/home.html", "../../view/components/navbar.html")
	if err != nil {
		log.Printf("Erro ao carregar template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// verifica se tem o header de autenticação
	// mesmo se tiver e for um token inválido, aparece as opções na navbar, mas não terá permissão para acessar

	// retorna erro se não encontrar
	_, err = r.Cookie("auth_token")
	IsAuthenticated := err == nil

	fmt.Println(IsAuthenticated)
	err = tmpl.Execute(w, homePageData{IsAuthenticated})
	if err != nil {
		log.Printf("Erro ao executar template: %v", err)
	}
}
