package auth

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

const (
	expirationMinutesNumber  = 10
	headerKey                = "Content-Type"
	contentType              = "application/json"
	unauthorizedMessage      = "авторизация не пройдена"
	missingAuthHeaderMessage = "отсутствует header авторизации"
)

func GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(expirationMinutesNumber * time.Minute)
	claims := &Claims{
		Login: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func AuthorizationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set(headerKey, contentType)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(missingAuthHeaderMessage))
			return
		}

		tokenString := authHeader[len("Bearer "):]
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			w.Header().Set(headerKey, contentType)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(unauthorizedMessage))
			return
		}

		next.ServeHTTP(w, r)
	}
}
