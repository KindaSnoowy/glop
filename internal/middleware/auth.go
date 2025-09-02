package authmiddleware

import (
	customErrors "blog_api/internal/errors"
	"blog_api/internal/repository"
	"context"
	"errors"
	"net/http"
	"strings"
)

type contextKey string

const UserIDKey contextKey = "user_id"

// inicializa o authmiddleware, injetando as dependências
func AuthMiddleware(sessionRepository *repository.SessionRepository) func(http.Handler) http.Handler {
	// recebe o handler http e retorna o handler modificado
	return func(next http.Handler) http.Handler {
		// cria uma função no handler
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "No auth header given", http.StatusUnauthorized)
				return
			}

			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				http.Error(w, "Auth header not formatted correctly", http.StatusUnauthorized)
				return
			}

			token := headerParts[1]

			userID, err := sessionRepository.IsTokenValid(token)
			if err != nil {
				if errors.Is(err, customErrors.ErrNotFound) {
					http.Error(w, "Invalid token", http.StatusUnauthorized)
					return
				} else if errors.Is(err, customErrors.ErrExpiredToken) {
					http.Error(w, "Expired token. Please log-in again.", http.StatusUnauthorized)
					return
				} else {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
