// Package middlewares -> Define os middlewares customizados do projeto.
package middlewares

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	customerrors "blog_api/internal/errors"
	"blog_api/internal/repository"
)

type contextKey string

const UserIDKey contextKey = "user_id"

// por enquanto, todos os usuários são administradores, por isso não existe verificação de permissões
// rota de criar usuário também não é pública

// AuthMiddleware -> Faz a autenticação, qualquer usuário com token válido passa
func AuthMiddleware(sessionRepository *repository.SessionRepository) func(http.Handler) http.Handler {
	// recebe o handler http e retorna o handler modificado
	return func(next http.Handler) http.Handler {
		// cria uma função no handler
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var token string

			log.Println("Request headers:", r.Header)
			log.Println("Authorization:", r.Header.Get("Authorization"))

			authHeader := r.Header.Get("Authorization")
			// verifica se um header Authorization foi passado
			if authHeader == "" {
				// se não foi passado, procura no cookie
				cookie, err := r.Cookie("session_token")
				if err != nil {
					http.Error(w, "No authorization cookie given", http.StatusUnauthorized)
					return
				}

				token = cookie.Value
			} else {
				// se foi passado, usa o token
				headerParts := strings.Split(authHeader, " ")
				if len(headerParts) != 2 || headerParts[0] != "Bearer" {
					http.Error(w, "Auth header not formatted correctly", http.StatusUnauthorized)
					return
				}

				token = headerParts[1]
			}

			userID, err := sessionRepository.IsTokenValid(token)
			if err != nil {
				if errors.Is(err, customerrors.ErrNotFound) {
					http.Error(w, "Invalid token", http.StatusUnauthorized)
					return
				} else if errors.Is(err, customerrors.ErrExpiredToken) {
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
