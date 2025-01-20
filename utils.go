package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type Error struct {
	Error string
}

func WriteJson(w http.ResponseWriter, status int, m any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(m)
}

func HashPassword(password string) (string, error) {
	byte, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(byte), err
}

func checkPassword(password string, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err
}

func checkLength(username string, password string) error {
	if len(username) < 8 || len(password) < 8 {
		return fmt.Errorf("not enougn long")
	}
	return nil
}
