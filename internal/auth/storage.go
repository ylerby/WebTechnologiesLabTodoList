package auth

import (
	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("secret_key")

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Claims struct {
	Login string `json:"login"`
	jwt.StandardClaims
}
