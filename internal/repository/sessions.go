package repository

import (
	customErrors "blog_api/internal/errors"
	"database/sql"
	"log"
	"time"
)

type SessionRepository struct {
	DB *sql.DB
}

func StartSessionRepository(db *sql.DB) (*SessionRepository, error) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS sessions (
						 token VARCHAR(255) NOT NULL PRIMARY KEY,
					    user_id INT NOT NULL,
					    created_at DATETIME NOT NULL,
					    expires_at DATETIME NOT NULL
					  );`

	if _, err := db.Exec(createTableSQL); err != nil {
		return nil, err
	}
	log.Println("Tabela 'sessions' pronta")

	return &SessionRepository{DB: db}, nil
}

func (s *SessionRepository) CreateSession(token string, userID int64, createdAt, expiresAt time.Time) error {
	_, err := s.DB.Exec(`INSERT INTO sessions (token, user_id, created_at, expires_at) VALUES (?, ?, ?, ?)`, token, userID, createdAt, expiresAt)
	return err
}

func (s *SessionRepository) DeleteSessionByID(id int64) error {
	result, err := s.DB.Exec(`DELETE FROM sessions WHERE id = ?`, id)
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

func (s *SessionRepository) DeleteSessionByUser(id int64) error {
	result, err := s.DB.Exec(`DELETE FROM sessions WHERE user_id = ?`, id)

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

func (s *SessionRepository) IsTokenValid(token string) (int64, error) {
	var expires_at time.Time
	var userID int64

	err := s.DB.QueryRow("SELECT user_id, expires_at FROM sessions WHERE token = ?", token).Scan(&userID, &expires_at)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, customErrors.ErrNotFound
		}
		return 0, err
	}

	if expires_at.Before(time.Now()) {
		// goroutine
		go func() {
			// limpa os tokens inválidos (poderia limpar somente o token que tentou ser utilizado)
			_, err := s.DB.Exec("DELETE FROM sessions WHERE expires_at < NOW();", token)
			if err != nil {
				log.Printf("Error while deleting expired token: %v", err)
			}
		}()
		return 0, customErrors.ErrExpiredToken
	}

	return userID, nil
}
