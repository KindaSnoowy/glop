package web

import (
	customErrors "blog_api/internal/errors"
	"blog_api/internal/models"
	"blog_api/internal/repository"
	passwords "blog_api/internal/utils"
	"errors"
	"log"
	"net/http"
	"text/template"
	"time"
)

const SESSION_DURATION_DAYS = 7

type LoginHandler struct {
	SessionRepository *repository.SessionRepository
	UserRepository    *repository.UserRepository
}

func StartLoginHandler(sessionRepository *repository.SessionRepository, userRepository *repository.UserRepository) *LoginHandler {
	return &LoginHandler{
		SessionRepository: sessionRepository,
		UserRepository:    userRepository,
	}
}

func (s *LoginHandler) GetLoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../../view/pages/login.html")
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

func (s *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusInternalServerError)
		return
	}

	loginDTO := models.LoginRequest{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	usuario, err := s.UserRepository.GetByUsername(loginDTO.Username)
	if err != nil {
		if errors.Is(err, customErrors.ErrNotFound) {
			http.Error(w, "Invalid user or password", http.StatusUnauthorized)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !passwords.CheckPasswordHash(loginDTO.Password, usuario.Password) {
		http.Error(w, "Invalid user or password", http.StatusUnauthorized)
		return
	}

	// senha válida, logou
	token, err := passwords.GenerateRandomToken(32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	expiresAt := time.Now().AddDate(0, 0, SESSION_DURATION_DAYS)
	err = s.SessionRepository.CreateSession(
		token,
		usuario.ID,
		time.Now(),
		expiresAt,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "session_token",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		Expires:  expiresAt,
	}
	http.SetCookie(w, &cookie)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
