package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	customerrors "blog_api/internal/errors"
	"blog_api/internal/models"
	"blog_api/internal/repository"

	"github.com/go-chi/chi/v5"
)

type PostHandlerAPI struct {
	Repository *repository.PostRepository
}

func StartPostHandler(repository *repository.PostRepository) *PostHandlerAPI {
	return &PostHandlerAPI{
		Repository: repository,
	}
}

func (s *PostHandlerAPI) CreatePost(w http.ResponseWriter, r *http.Request) {
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

	id, err := s.Repository.Create(&postDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	postDTO.ID = id

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(postDTO)
}

func (s *PostHandlerAPI) GetPosts(w http.ResponseWriter, r *http.Request) {
	var postFilter models.PostFilters
	if err := json.NewDecoder(r.Body).Decode(&postFilter); err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}
	fmt.Printf("%+v\n", postFilter)

	postsDatabase, err := s.Repository.GetAll(&postFilter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(postsDatabase); err != nil {
		http.Error(w, "Error when serializing response", http.StatusInternalServerError)
	}
}

func (s *PostHandlerAPI) GetPostById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	post, err := s.Repository.GetByID(id)
	if err != nil {
		if errors.Is(err, customerrors.ErrNotFound) {
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

func (s *PostHandlerAPI) Update(w http.ResponseWriter, r *http.Request) {
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
	postInterno.UpdatedAt = time.Now()

	err = s.Repository.Update(id, postInterno)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if err == customerrors.ErrNotFound {
		http.Error(w, "Post with that ID was not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(postInterno)
}

func (s *PostHandlerAPI) DeletePost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = s.Repository.Delete(id)
	if err == customerrors.ErrNotFound {
		http.Error(w, "Post with that ID was not found", http.StatusNotFound)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
