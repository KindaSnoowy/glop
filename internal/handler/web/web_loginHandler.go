package web

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	customerrors "blog_api/internal/errors"
	"blog_api/internal/models"
	"blog_api/internal/services"
)

type LoginHandler struct {
	AuthService *services.AuthService
}

func StartLoginHandler(authService *services.AuthService) *LoginHandler {
	return &LoginHandler{
		AuthService: authService,
	}
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

func (s *LoginHandler) WebLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	if username == "" || password == "" {
		http.Error(w, "empty inputs", http.StatusBadRequest)
		return
	}

	loginResponse, err := s.AuthService.AuthenticateUser(
		models.LoginRequest{
			Username: username, Password: password,
		},
	)
	if err != nil {
		if err == customerrors.ErrInvalidToken || err == customerrors.ErrNotFound {
			http.Error(w, "Username or password are invalid", http.StatusUnauthorized)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("web login, username = %s ; password = %s", username, password)
	http.SetCookie(w, &http.Cookie{
		Name:  "auth_token",
		Value: loginResponse.Token,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(loginResponse)
}
