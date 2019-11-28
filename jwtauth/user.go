package jwtauth

import "github.com/dgrijalva/jwt-go"

type Credentials struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type User struct {
	id        string
	email     string
	password  string
	isDeleted int
}

type Response struct {
	Data string `json:"data"`
}

type Token struct {
	Token string `json:"token"`
}

type UserInfo struct {
	id    string
	email string
}

type CustomClaims struct {
	id    string `json:"id"`
	email string `json:"mail"`
	jwt.StandardClaims
}
