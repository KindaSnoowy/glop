package handlers

import (
	"blog_api/internal/models"
	"blog_api/internal/repository"
	"encoding/json"
	"net/http"
	"time"
)

type PostHandler struct {
	Repository *repository.PostRepository
}

func StartPostHandler(repository *repository.PostRepository) *PostHandler {
	return &PostHandler{
		Repository: repository,
	}
}

func (s *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var postCreateDTO models.PostCreateDTO
	if err := json.NewDecoder(r.Body).Decode(&postCreateDTO); err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	postDTO := models.Post{
		ID:        0,
		Title:     postCreateDTO.Title,
		Content:   postCreateDTO.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	post, err := s.Repository.Create(postDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func (s *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	postsDatabase, err := s.Repository.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(postsDatabase); err != nil {
		http.Error(w, "Error when serializing response", http.StatusInternalServerError)
	}
}
