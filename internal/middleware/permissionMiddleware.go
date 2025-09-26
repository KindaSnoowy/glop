package middlewares

import (
	customErrors "blog_api/internal/errors"
	"blog_api/internal/repository"
	"net/http"
	"strconv"
)

func PermissionMiddleware(userRepository *repository.UserRepository) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value(UserIDKey).(string)

			userID_asInteger, err := strconv.Atoi(userID)
			if err != nil {
				http.Error(w, "Error while converting userID from context", http.StatusInternalServerError)
				return
			}

			user, err := userRepository.GetByID((userID_asInteger))
			if err != nil {
				if err == customErrors.ErrNotFound {
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
