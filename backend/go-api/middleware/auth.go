package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if tokenStr == "" {
			http.Error(w, "Unauthorized", 401)
			return
		}

		_, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte("replace-with-env-secret"), nil
		})

		if err != nil {
			http.Error(w, "Invalid token", 401)
			return
		}

		next.ServeHTTP(w, r)
	})
}
