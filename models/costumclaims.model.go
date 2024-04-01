package models

import "github.com/golang-jwt/jwt/v5"

type JwtCustomClaims struct {
	ID int `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role int `json:"role"`
	jwt.RegisteredClaims
}