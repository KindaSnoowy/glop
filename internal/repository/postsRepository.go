// Package repository -> repositórios do projeto, responsável somente pelas requisições ao banco.
package repository

import (
	"database/sql"
	"log"

	customerrors "blog_api/internal/errors"
	"blog_api/internal/models"
)

type PostRepository struct {
	DB *sql.DB
}

func StartPostRepository(db *sql.DB) (*PostRepository, error) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS posts (
				"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, "title" TEXT, "content" TEXT,
				"createdAt" DATETIME,"updatedAt" DATETIME
			);`
	if _, err := db.Exec(createTableSQL); err != nil {
		return nil, err
	}
	log.Println("Tabela 'posts' pronta")

	return &PostRepository{
		DB: db,
	}, nil
}

// Create -> retorna id ao invés de ponteiro para simplificar a assinatura
// recebe ponteiro pois só lê o objeto (não edita ele em nenhum momento)
func (s *PostRepository) Create(post *models.Post) (int64, error) {
	result, err := s.DB.Exec(`INSERT INTO posts (title, content, createdAt, updatedAt) VALUES (?, ?, ?, ?)`, post.Title, post.Content, post.CreatedAt, post.UpdatedAt)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *PostRepository) GetByID(id int) (*models.Post, error) {
	row := s.DB.QueryRow(`SELECT id, title, content, createdAt, updatedAt FROM posts WHERE id = ?`, id)

	var post models.Post
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, customerrors.ErrNotFound
		}

		return nil, err
	}

	return &post, nil
}

// GetAll busca posts com filtros, paginação e busca.
func (s *PostRepository) GetAll(filters *models.PostFilters) ([]models.Post, error) {
	var args []any
	var query string

	if filters.ShortContent {
		query = `SELECT id, title, SUBSTRING(content, 1, 500) as content, createdAt, updatedAt FROM posts`
	} else {
		query = `SELECT id, title, content, createdAt, updatedAt FROM posts`
	}

	query += " WHERE 1=1"
	if filters.Search != "" {
		query += " AND (title LIKE ? OR content LIKE ?)"
		likeTerm := "%" + filters.Search + "%"
		args = append(args, likeTerm, likeTerm)
	}

	var orderCreated string
	if filters.OrderCreated {
		orderCreated = "ASC"
	} else {
		orderCreated = "DESC"
	}
	query += " ORDER BY createdAt " + orderCreated

	query += " LIMIT ? OFFSET ?"
	offset := max((filters.Page-1)*filters.Limit, 0)
	args = append(args, filters.Limit, offset)

	rows, err := s.DB.Query(query, args...)
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

func (s *PostRepository) Update(id int, postDTO *models.Post) error {
	result, err := s.DB.Exec(`UPDATE posts SET title = ?, content = ?, updatedAt = ? WHERE id = ?`, postDTO.Title, postDTO.Content, postDTO.UpdatedAt, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return customerrors.ErrNotFound
	}

	return nil
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
		return customerrors.ErrNotFound
	}

	return nil
}
