package api

import (
	"encoding/json"
	"net/http"

	customerrors "blog_api/internal/errors"
	"blog_api/internal/models"
	"blog_api/internal/services"
)

const SessionDurationDays = 7

type LoginHandlerAPI struct {
	AuthService *services.AuthService
}

func StartLoginHandler(authService *services.AuthService) *LoginHandlerAPI {
	return &LoginHandlerAPI{
		AuthService: authService,
	}
}

func (s *LoginHandlerAPI) Login(w http.ResponseWriter, r *http.Request) {
	var loginDTO models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginDTO); err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	loginResponse, err := s.AuthService.AuthenticateUser(loginDTO)
	if err == customerrors.ErrInvalidToken || err == customerrors.ErrNotFound {
		// se erro == invalid token, a senha está errada
		// se erro == não encontrado, o usuário está errado
		// retorna o mesmo erro pros dois por questões de segurança

		http.Error(w, "Username or password are invalid", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(loginResponse)
}
