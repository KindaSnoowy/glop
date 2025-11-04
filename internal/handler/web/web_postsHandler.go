package web

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	customerrors "blog_api/internal/errors"
	"blog_api/internal/models"
	"blog_api/internal/repository"

	"github.com/go-chi/chi/v5"
)

type PostHandler struct {
	Repository *repository.PostRepository
}

// render datas
type PostsPageData struct {
	IsAuthenticated bool
	Posts           []models.Post
	NextPage        int
}

type PostPageData struct {
	IsAuthenticated bool
	Post            *models.Post
}

func StartPostHandler(repository *repository.PostRepository) *PostHandler {
	return &PostHandler{
		Repository: repository,
	}
}

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

	_, err = r.Cookie("auth_token")
	IsAuthenticated := err == nil

	data := PostsPageData{
		IsAuthenticated: IsAuthenticated,
		Posts:           posts,
		NextPage:        nextPage,
	}

	// se for requisição do htmx, renderiza só o componente da lista de posts
	if r.Header.Get("HX-Request") == "true" {
		tmpl, err := template.ParseFiles("../../view/components/post_list.html")
		if err != nil {
			log.Printf("Erro ao carregar template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		fmt.Printf("%+v\n", data)
		err = tmpl.Execute(w, data)
		if err != nil {
			log.Printf("Erro ao executar template: %v", err)
		}
		return
	}

	tmpl, err := template.ParseFiles("../../view/pages/posts.html",
		"../../view/components/post_list.html", "../../view/components/navbar.html")
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

func (s *PostHandler) GetPostIDPage(
	w http.ResponseWriter, r *http.Request,
) {
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

	tmpl, err := template.ParseFiles("../../view/pages/post.html", "../../view/components/navbar.html")
	if err != nil {
		log.Printf("Erro ao carregar template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	_, err = r.Cookie("auth_token")
	IsAuthenticated := err == nil

	data := PostPageData{
		IsAuthenticated: IsAuthenticated,
		Post:            post,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Erro ao executar template: %v", err)
	}
}

func (s *PostHandler) WebCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	title := r.Form.Get("title")
	content := r.Form.Get("content")

	post := models.Post{Title: title, Content: content}
	fmt.Println("conteudo:", title, content)
	_, err = s.Repository.Create(&post)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("HX-Redirect", "/posts/")
	w.WriteHeader(http.StatusOK)
}

func (s *PostHandler) WebUpdatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	title := r.Form.Get("title")
	content := r.Form.Get("content")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	postInterno, err := s.Repository.GetByID(id)
	if err != nil {
		if err == customerrors.ErrNotFound {
			http.Error(w, "Post with ID not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	postDTO := models.Post{
		ID:        postInterno.ID,
		Title:     title,
		Content:   content,
		CreatedAt: postInterno.CreatedAt,
		UpdatedAt: time.Now(),
	}

	err = s.Repository.Update(id, &postDTO)
	if err != nil {
		if err == customerrors.ErrNotFound {
			http.Error(w, "Post with ID not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("HX-Redirect", "/posts/")
	w.WriteHeader(http.StatusOK)
}
