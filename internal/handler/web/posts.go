package web

import (
	"blog_api/internal/models"
	"blog_api/internal/repository"
	"log"
	"net/http"
	"strconv"
	"text/template"

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

// all posts
func (s *PostHandler) GetPostsPage(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	posts, err := s.Repository.GetAll(
		&models.PostFilters{
			ShortContent: true,
			Limit:        5,
			Page:         page,
		},
	)

	if err != nil {
		log.Printf("Erro ao buscar posts: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var nextPage int
	if len(posts) == 5 {
		nextPage = page + 1
	} // will be 0 if theres not a next page

	data := models.PostPageData{
		Posts:    posts,
		NextPage: nextPage,
	}

	if r.Header.Get("HX-Request") == "true" {
		tmpl, err := template.ParseFiles("../../view/components/post_list.html")
		if err != nil {
			log.Printf("Erro ao carregar template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			log.Printf("Erro ao executar template: %v", err)
		}
		return
	}

	tmpl, err := template.ParseFiles("../../view/pages/posts.html")
	if err != nil {
		log.Printf("Erro ao carregar template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Erro ao executar template: %v", err)
	}
}

// post by id
func (s *PostHandler) GetPostIdPage(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	post, err := s.Repository.GetByID(id)
	if post == nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Erro ao carregar post: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("../../view/pages/post.html")
	if err != nil {
		log.Printf("Erro ao carregar template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, post)
	if err != nil {
		log.Printf("Erro ao executar template: %v", err)
	}
}
