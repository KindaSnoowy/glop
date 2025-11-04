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

type loginPageData struct {
	Error string
}

func StartLoginHandler(authService *services.AuthService) *LoginHandler {
	return &LoginHandler{
		AuthService: authService,
	}
}

func (s *LoginHandler) RenderLoginComponent(w http.ResponseWriter, errorMessage string) {
	tmpl, err := template.ParseFiles("../../view/components/login_form.html")
	if err != nil {
		log.Printf("Erro ao carregar página de login: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, loginPageData{Error: errorMessage})
	if err != nil {
		log.Printf("Erro ao executar template: %v", err)
	}
}

func (s *LoginHandler) GetLoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../../view/pages/login.html", "../../view/components/login_form.html")
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
		s.RenderLoginComponent(w, "Usuário ou Senha nulos")
		return
	}

	loginResponse, err := s.AuthService.AuthenticateUser(
		models.LoginRequest{
			Username: username, Password: password,
		},
	)
	if err != nil {
		if err == customerrors.ErrInvalidToken || err == customerrors.ErrNotFound {
			s.RenderLoginComponent(w, "Usuário ou Senha inválidos")
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
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("<title>Sucesso!</title><h1>Login realizado com sucesso!</h1>")
}

func (s *LoginHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
