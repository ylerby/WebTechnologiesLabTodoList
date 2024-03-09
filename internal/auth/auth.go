package auth

import (
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
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
			_, err := w.Write([]byte(missingAuthHeaderMessage))
			if err != nil {
				log.Printf("ошибка при получении ответа - %s", err)
				return
			}

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
			_, err = w.Write([]byte(unauthorizedMessage))
			if err != nil {
				log.Printf("ошибка при получении ответа - %s", err)
				return
			}

			return
		}

		next.ServeHTTP(w, r)
	}
}
