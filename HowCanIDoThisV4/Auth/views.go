package auth

import "github.com/golang-jwt/jwt"

type Claims struct {
	jwt.StandardClaims
}

type User struct {
	Id       uint32
	Password string
	Email    string
}
