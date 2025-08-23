package repository

import (
	"blog_api/internal"
	"blog_api/internal/models"
	"database/sql"
	"errors"
	"time"
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

func (s *PostRepository) GetByID(id int) (*models.Post, error) {
	row := s.DB.QueryRow(`SELECT id, title, content, createdAt, updatedAt FROM posts WHERE id = ?`, id)

	var post models.Post
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, internal.ErrNotFound
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

func (s *PostRepository) Update(id int, postDTO models.Post) (*models.Post, error) {
	postDTO.UpdatedAt = time.Now()
	result, err := s.DB.Exec(`UPDATE posts SET title = ?, content = ?, updatedAt = ? WHERE id = ?`, postDTO.Title, postDTO.Content, postDTO.UpdatedAt, id)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, nil
	}

	return &postDTO, nil
}

func (s *PostRepository) Delete(id int) error {
	result, err := s.DB.Exec(`DELETE FROM posts WHERE id = ?`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("Not Found")
	}

	return nil
}
