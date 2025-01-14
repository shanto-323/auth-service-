package main

import (
	"fmt"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	fmt.Println("env working")
	if err != nil {
		fmt.Println("env not working %v", err)
		panic(err)
	}
}

func createJWT(username string, hashed_password string) (string, error) {
	claims := &jwt.MapClaims{
		"exp":      time.Now().Add(time.Hour).Unix(),
		"username": username,
		"password": hashed_password,
	}

	mySigningKey := os.Getenv("JWTKEY")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(mySigningKey))
}

func createRefreshToken(auth *Auth) (string, error) {
	claims := &jwt.MapClaims{
		"id":       auth.Id,
		"username": auth.Username,
		"password": auth.Password.HashedPassword,
		"iat":      time.Now().Unix(),
	}

	mySigningKey := os.Getenv("JWTKEY")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(mySigningKey))
}
