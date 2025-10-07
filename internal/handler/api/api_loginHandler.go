package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	customerrors "blog_api/internal/errors"
	"blog_api/internal/models"
	"blog_api/internal/repository"
	passwords "blog_api/internal/utils"
)

const SessionDurationDays = 7

type LoginHandlerAPI struct {
	SessionRepository *repository.SessionRepository
	UserRepository    *repository.UserRepository
}

func StartLoginHandler(sessionRepository *repository.SessionRepository, userRepository *repository.UserRepository) *LoginHandlerAPI {
	return &LoginHandlerAPI{
		SessionRepository: sessionRepository,
		UserRepository:    userRepository,
	}
}

func (s *LoginHandlerAPI) Login(w http.ResponseWriter, r *http.Request) {
	var loginDTO models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginDTO); err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	usuario, err := s.UserRepository.GetByUsername(loginDTO.Username)
	if err != nil {
		if errors.Is(err, customerrors.ErrNotFound) {
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

	err = s.SessionRepository.CreateSession(
		token,
		usuario.ID,
		time.Now(),
		time.Now().AddDate(0, 0, SessionDurationDays),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	loginResponse := models.LoginResponse{
		Token: token,
		User: models.UserResponseDTO{
			ID:       usuario.ID,
			Name:     usuario.Name,
			Username: usuario.Username,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(loginResponse)
}
