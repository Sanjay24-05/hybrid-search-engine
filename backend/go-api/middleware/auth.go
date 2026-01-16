package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		if tokenStr == "" || tokenStr == authHeader {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims := jwt.MapClaims{}

		_, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// ✅ Extract user identifier (using email for now)
		userID, ok := claims["email"].(string)
		if !ok || userID == "" {
			http.Error(w, "Invalid token payload", http.StatusUnauthorized)
			return
		}

		// ✅ Inject user ID for downstream handlers (upload/search)
		r.Header.Set("X-User-ID", userID)

		next.ServeHTTP(w, r)
	})
}
