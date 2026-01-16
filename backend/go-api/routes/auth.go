package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("replace-with-env-secret")

func Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string
		Password string
	}
	json.NewDecoder(r.Body).Decode(&creds)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("password"), 10)
	if bcrypt.CompareHashAndPassword(hashed, []byte(creds.Password)) != nil {
		http.Error(w, "Invalid credentials", 401)
		return
	}

	claims := jwt.MapClaims{
		"email": creds.Email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString(jwtKey)

	json.NewEncoder(w).Encode(map[string]string{"token": tokenStr})
}
