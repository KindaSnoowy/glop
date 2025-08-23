package handlers

import (
	"blog_api/internal"
	"blog_api/internal/models"
	"blog_api/internal/repository"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
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

func (s *PostHandler) GetPostById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	post, err := s.Repository.GetByID(id)
	if err != nil {
		if errors.Is(err, internal.ErrNotFound) {
			http.Error(w, "Post with that ID was not found", http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(post)
}

func (s *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var postDTO models.PostUpdateDTO
	if err := json.NewDecoder(r.Body).Decode(&postDTO); err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	postInterno, err := s.Repository.GetByID(id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if postInterno == nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	postInterno.Title = postDTO.Title
	postInterno.Content = postDTO.Content

	post, err := s.Repository.Update(id, *postInterno)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(post)
}

func (s *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = s.Repository.Delete(id)
	if err == internal.ErrNotFound {
		http.Error(w, "Post with that ID was not found", http.StatusNotFound)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
