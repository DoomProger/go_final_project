package main

// import (
// 	"crypto/sha256"
// 	"encoding/hex"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"os"
// 	"time"

// 	"github.com/go-chi/jwtauth/v5"
// 	"github.com/go-chi/render"
// 	"github.com/golang-jwt/jwt"
// )

// var tokenAuth *jwtauth.JWTAuth
// var todoPasswordHash string
// var todoPassword = os.Getenv("TODO_PASSWORD")

// func init() {
// 	// Fetch and hash the password from environment variable
// 	// todoPassword := os.Getenv("TODO_PASSWORD")
// 	if todoPassword == "" {
// 		fmt.Println("TODO_PASSWORD environment variable is not set")
// 		os.Exit(1)
// 	}
// 	todoPasswordHash = hashPassword(todoPassword)

// 	// Create a new JWTAuth instance
// 	tokenAuth = jwtauth.New("HS256", []byte("123"), nil)
// }

// func verifyPassword(input, expected string) bool {
// 	return input == expected
// }

// func computeChecksum(str string) string {
// 	hash := sha256.Sum256([]byte(str))
// 	return hex.EncodeToString(hash[:])
// }

// type Credentials struct {
// 	Password string `json:"password"`
// }

// type ResponseAuth struct {
// 	Token string `json:"token,omitempty"`
// 	Error string `json:"error,omitempty"`
// }

// func handleSignIn(w http.ResponseWriter, r *http.Request) {
// 	var creds Credentials
// 	err := json.NewDecoder(r.Body).Decode(&creds)
// 	if err != nil || creds.Password == "" {
// 		writeJSONError(w, http.StatusUnauthorized, "Invalid request")
// 		return
// 	}

// 	if hashPassword(creds.Password) != todoPasswordHash {
// 		writeJSONError(w, http.StatusUnauthorized, "Invalid password")
// 		return
// 	}

// 	claims := map[string]interface{}{"checksum": todoPasswordHash}
// 	_, tokenString, _ := tokenAuth.Encode(claims)

// 	// json.NewEncoder(w).Encode(ResponseAuth{Token: tokenString})
// 	// rt := ResponseAuth{Token: tokenString}
// 	// log.Println("token:", rt)
// 	if verifyPassword(creds.Password, todoPassword) {
// 		_, tokenString, _ := tokenAuth.Encode(jwt.MapClaims{
// 			"checksum": computeChecksum(todoPassword),
// 			"exp":      time.Now().Add(15 * time.Minute).Unix(),
// 		})
// 		render.JSON(w, r, map[string]string{"token": tokenString})
// 	} else {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 	}

// 	render.JSON(w, r, ResponseAuth{Token: tokenString})
// }

// func hashPassword(password string) string {
// 	hasher := sha256.New()
// 	hasher.Write([]byte(password))
// 	return hex.EncodeToString(hasher.Sum(nil))
// }
//
