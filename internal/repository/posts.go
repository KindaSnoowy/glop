package repository

import (
	"blog_api/internal/models"
	"database/sql"
)

type PostRepository struct {
	DB *sql.DB
}

func StartPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{
		DB: db,
	}
}

func (s *PostRepository) Create(post models.Post) (*models.Post, error) {
	// post é passado como valor, não como ponteiro, pra criar uma cópia

	result, err := s.DB.Exec(`INSERT INTO posts (title, content, createdAt, updatedAt) VALUES (?, ?, ?, ?)`, post.Title, post.Content, post.CreatedAt, post.UpdatedAt)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	post.ID = int(id)

	return &post, nil
}

func (s *PostRepository) GetByID(id int32) (*models.Post, error) {
	row := s.DB.QueryRow(`SELECT id, title, content, createdAt, updatedAt FROM posts WHERE id = ?`, id)

	var post models.Post
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &post, nil
}

func (s *PostRepository) GetAll() ([]models.Post, error) {
	rows, err := s.DB.Query(`SELECT id, title, content, createdAt, updatedAt FROM posts ORDER BY createdAt DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []models.Post{}
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}
