package models

import "github.com/golang-jwt/jwt/v4"

type User struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type JWTCustomClaims struct {
    Username string `json:"username"`
    jwt.RegisteredClaims
}
