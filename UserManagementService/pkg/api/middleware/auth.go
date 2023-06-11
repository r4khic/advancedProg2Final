package middleware

import (
	"advancedProg2Final/UserManagementService/pkg/service"
	"net/http"
)

func Auth(s *service.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")

			if token == "" {
				http.Error(w, "Authorization token must be present", http.StatusUnauthorized)
				return
			}

			_, err := s.ValidateToken(token)
			if err != nil {
				http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
