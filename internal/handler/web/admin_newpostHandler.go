package web

import (
	customerrors "blog_api/internal/errors"
	"blog_api/internal/models"
	"blog_api/internal/repository"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/go-chi/chi/v5"
)

type NewPostHandler struct {
	PostRepository *repository.PostRepository
}

type newPostPageData struct {
	IsAuthenticated bool
	IsEditing       bool
	Post            models.Post
}

func StartNewPostHandler(postRepository *repository.PostRepository) *NewPostHandler {
	return &NewPostHandler{PostRepository: postRepository}
}

func (s *NewPostHandler) GetNewPostPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../../view/pages/newpost.html", "../../view/components/navbar.html")
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

	// verifica se passou parametro na url (edit/{id}), se passou está editando, se não, não está!
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	isEditing := err == nil

	// carrega o post se estiver editando, post vai ser nulo se não estiver
	var post models.Post
	if isEditing {
		postFound, err := s.PostRepository.GetByID(id)
		if err != nil {
			if err == customerrors.ErrNotFound {
				http.Error(w, "Post not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		post = *postFound
	}

	renderData := newPostPageData{
		IsAuthenticated: IsAuthenticated,
		IsEditing:       isEditing,
		Post:            post,
	}

	err = tmpl.Execute(w, renderData)
	if err != nil {
		log.Printf("Erro ao executar template: %v", err)
	}
}
