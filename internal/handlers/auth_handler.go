package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"gofinalproject/config"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var todoPasswordHash string
var todoPassword = os.Getenv("TODO_PASSWORD")

var jwtKey = []byte("SecretJWTKey_2142")

type Credentials struct {
	Password string `json:"password"`
}

type Claims struct {
	PasswordHash string `json:"password_hash"`
	jwt.RegisteredClaims
}

func computeChecksum(str string) string {
	hash := sha256.Sum256([]byte(str))
	return hex.EncodeToString(hash[:])
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
	todoPasswordHash = hashPassword(todoPassword)
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	// func (th *TaskHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var creds Credentials

	log.Println("SignIn hand")

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Check the password
	if hashPassword(creds.Password) != todoPasswordHash {
		writeJSONError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	// Create the JWT token
	expirationTime := time.Now().Add(config.TokenTTL * time.Hour)
	claims := &Claims{
		PasswordHash: todoPasswordHash,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	// http.SetCookie(w, &http.Cookie{
	// 	Name:     "token",
	// 	Value:    tokenStr,
	// 	Expires:  expirationTime,
	// 	HttpOnly: true,
	// })

	log.Println(tokenStr, "<-- token")

	// http.Redirect(w, r, "/", http.StatusSeeOther)

	// json.NewEncoder(w).Encode(creds)
	json.NewEncoder(w).Encode(map[string]string{"token": tokenStr})
}

// func (th *TaskHandler) AuthMiddleware(next http.Handler) http.Handler {
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Println("AuthMiddleware hand")

		tokenCookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenStr := tokenCookie.Value
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		log.Println(tokenStr)

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if claims.PasswordHash != todoPasswordHash {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
