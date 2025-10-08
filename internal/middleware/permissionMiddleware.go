// Package middlewares -> Define os middlewares customizados do projeto.
package middlewares

import (
	"net/http"

	customerrors "blog_api/internal/errors"
	"blog_api/internal/repository"
)

func PermissionMiddleware(userRepository *repository.UserRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userIDValue := r.Context().Value(UserIDKey)
			userID, ok := userIDValue.(int)
			if !ok {
				http.Error(w, "Invalid userID type in context", http.StatusInternalServerError)

				return
			}

			user, err := userRepository.GetByID(userID)
			if err != nil {
				if err == customerrors.ErrNotFound {
					http.Error(w, "User not found.", http.StatusInternalServerError)
					return
				} else {
					http.Error(w, "Error while getting user", http.StatusInternalServerError)
					return
				}
			}

			if user.IsAdmin {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "User do not have permission to access this route", http.StatusUnauthorized)
				return
			}
		})
	}
}
