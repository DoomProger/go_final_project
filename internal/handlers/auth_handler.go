package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var todoPassword = os.Getenv("TODO_PASSWORD")

var jwtKey []byte

type Credentials struct {
	Password string `json:"password"`
}

func hashPassword(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

func init() {
	if todoPassword == "" {
		log.Println("TODO_PASSWORD environment variable is not set")
		os.Exit(1)
	}
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	var creds Credentials

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if creds.Password != todoPassword {
		writeJSONError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	jwtKey = []byte(hashPassword(creds.Password))

	jwtToken := jwt.New(jwt.SigningMethodHS256)

	tokenStr, err := jwtToken.SignedString(jwtKey)

	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": tokenStr})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenCookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenStr := tokenCookie.Value
		jwtToken, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			http.Error(w, "Cannot parse token", http.StatusUnauthorized)
			return
		}

		if !jwtToken.Valid {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
