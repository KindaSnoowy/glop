package repository

import (
	customErrors "blog_api/internal/errors"
	"blog_api/internal/models"
	"database/sql"
	"log"
)

type UserRepository struct {
	DB *sql.DB
}

func StartUserRepository(db *sql.DB) (*UserRepository, error) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
				"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
				"name" VARCHAR(255) NOT NULL, 
				"username" VARCHAR(50) NOT NULL UNIQUE,
				"password" VARCHAR(255) NOT NULL
			);`

	if _, err := db.Exec(createTableSQL); err != nil {
		return nil, err
	}
	log.Println("Tabela 'users' pronta")

	return &UserRepository{DB: db}, nil
}

func (s *UserRepository) Create(user *models.User) (int64, error) {
	result, err := s.DB.Exec(`INSERT INTO users (name, username, password) VALUES (?, ?, ?)`, user.Name, user.Username, user.Password)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *UserRepository) GetByID(id int) (*models.User, error) {
	row := s.DB.QueryRow(`SELECT id, name, username, password FROM users WHERE id = ?`, id)

	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, customErrors.ErrNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (s *UserRepository) GetByUsername(username string) (*models.User, error) {
	row := s.DB.QueryRow(`SELECT id, name, username, password FROM users WHERE username = ?`, username)

	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, customErrors.ErrNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (s *UserRepository) Update(id int, userDTO *models.User) error {
	result, err := s.DB.Exec(`UPDATE users SET name = ?, username = ?, password = ? WHERE id = ?`, userDTO.Name, userDTO.Username, userDTO.Password, id)
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

func (s *UserRepository) Delete(id int) error {
	result, err := s.DB.Exec(`DELETE FROM users WHERE id = ?`, id)
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
