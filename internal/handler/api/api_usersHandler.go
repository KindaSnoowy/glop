package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	customerrors "blog_api/internal/errors"
	authmiddleware "blog_api/internal/middleware"
	"blog_api/internal/models"
	"blog_api/internal/repository"
	passwords "blog_api/internal/utils"

	"github.com/go-chi/chi/v5"
)

type UserHandlerAPI struct {
	Repository *repository.UserRepository
}

func StartUserHandler(repository *repository.UserRepository) *UserHandlerAPI {
	return &UserHandlerAPI{
		Repository: repository,
	}
}

func (s *UserHandlerAPI) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	user, err := s.Repository.GetByID(id)
	if err != nil {
		if errors.Is(err, customerrors.ErrNotFound) {
			http.Error(w, "User with that ID was not found", http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	userResponse := models.UserResponseDTO{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
	}
	json.NewEncoder(w).Encode(userResponse)
}

func (s *UserHandlerAPI) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	user, err := s.Repository.GetByUsername(username)
	if err != nil {
		if errors.Is(err, customerrors.ErrNotFound) {
			http.Error(w, "User with that Username was not found", http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	userResponse := models.UserResponseDTO{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
	}
	json.NewEncoder(w).Encode(userResponse)
}

func (s *UserHandlerAPI) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userCreateDTO models.UserCreateDTO
	err := json.NewDecoder(r.Body).Decode(&userCreateDTO)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	encryptedPassword, err := passwords.HashPassword(userCreateDTO.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userDTO := models.User{
		ID:       0,
		Name:     userCreateDTO.Name,
		Username: userCreateDTO.Username,
		Password: encryptedPassword,
	}

	id, err := s.Repository.Create(&userDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userResponse := models.UserResponseDTO{
		ID:       id,
		Name:     userDTO.Name,
		Username: userDTO.Username,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userResponse)
}

func (s *UserHandlerAPI) UpdateUser(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(authmiddleware.UserIDKey).(int64)
	if !ok {
		http.Error(w, "Não foi possível identificar o usuário logado", http.StatusInternalServerError)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var userDTO models.UserUpdateDTO
	if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	userInterno, err := s.Repository.GetByID(id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if userInterno == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if passwords.CheckPasswordHash(userDTO.Password, userInterno.Password) {
		http.Error(w, "New password cannot be the same as the old one", http.StatusNotFound)
		return
	}

	hash, err := passwords.HashPassword(userDTO.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userInterno.Name = userDTO.Name
	userInterno.Username = userDTO.Username
	userInterno.Password = hash

	err = s.Repository.Update(id, userInterno)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if err == customerrors.ErrNotFound {
		http.Error(w, "User with that ID was not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userInterno)
}
