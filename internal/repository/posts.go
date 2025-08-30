package repository

import (
	customErrors "blog_api/internal/errors"
	"blog_api/internal/models"
	"database/sql"
	"log"
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

// retorna id ao invés de ponteiro para simplificar a assinatura
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
			return nil, customErrors.ErrNotFound
		}

		return nil, err
	}

	return &post, nil
}

func (s *PostRepository) GetAll(filters *models.PostFilters) ([]models.Post, error) {
	var query string
	if filters.ShortContent {
		query = `SELECT id, title, SUBSTRING(content, 1, 500) as content, createdAt, updatedAt FROM posts ORDER BY createdAt DESC LIMIT ? OFFSET ?`
	} else {
		query = `SELECT id, title, content, createdAt, updatedAt FROM posts ORDER BY createdAt DESC LIMIT ? OFFSET ?`
	}

	rows, err := s.DB.Query(query, filters.Limit, filters.Page*filters.Limit)
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
		return customErrors.ErrNotFound
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
		return customErrors.ErrNotFound
	}

	return nil
}
