// Package services ->
package services

import (
	customerrors "blog_api/internal/errors"
	"blog_api/internal/models"
	"blog_api/internal/repository"
	passwords "blog_api/internal/utils"
	"time"
)

type AuthService struct {
	SessionRepository *repository.SessionRepository
	UserRepository    *repository.UserRepository
}

const SessionDurationDays = 7

func StartAuthService(sessionRepository *repository.SessionRepository,
	userRepository *repository.UserRepository) *AuthService {
	return &AuthService{
		SessionRepository: sessionRepository,
		UserRepository:    userRepository,
	}
}

func (s *AuthService) AuthenticateUser(loginRequest models.LoginRequest) (*models.LoginResponse, error) {
	usuario, err := s.UserRepository.GetByUsername(loginRequest.Username)
	if err != nil {
		return nil, customerrors.ErrNotFound
	}

	if !passwords.CheckPasswordHash(loginRequest.Password, usuario.Password) {
		return nil, customerrors.ErrInvalidToken
	}

	// senha válida, logou
	token, err := passwords.GenerateRandomToken(32)
	if err != nil {
		return nil, err
	}

	err = s.SessionRepository.CreateSession(
		token,
		usuario.ID,
		time.Now(),
		time.Now().AddDate(0, 0, SessionDurationDays),
	)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Token: token,
		User: models.UserResponseDTO{
			ID:       usuario.ID,
			Name:     usuario.Name,
			Username: usuario.Username,
			IsAdmin:  usuario.IsAdmin,
		},
	}, nil
}
