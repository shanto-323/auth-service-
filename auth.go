package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	fmt.Println("env working")
	if err != nil {
		fmt.Printf("env not working %v", err)
		panic(err)
	}
}

func createJWT(auth *Auth) (string, error) {
	claims := &jwt.MapClaims{
		//"exp":      time.Now().Add(time.Hour).Unix(),
		"id":         auth.Id,
		"username":   auth.Username,
		"password":   auth.Password.HashedPassword,
		"created_at": time.Now(),
	}

	mySigningKey := os.Getenv("JWTKEY")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(mySigningKey))
}

func withJwt(handlerfunc http.HandlerFunc, storage Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("x-jwt-token")
		token, err := validateJwt(tokenString)
		if err != nil {
			WriteJson(w, http.StatusForbidden, Error{Error: err.Error()})
			return
		}

		idString := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			WriteJson(w, http.StatusBadRequest, err.Error())
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		claimsUser := &Auth{
			Id:       int64((claims["id"]).(float64)),
			Username: claims["username"].(string),
			Password: Password{
				HashedPassword: claims["password"].(string),
			},
		}

		if claimsUser.Id != int64(id) {
			WriteJson(w, http.StatusBadRequest, Error{Error: "id not matched"})
			return
		}

		databaseUser, err := storage.GetAccountById(id)
		if err != nil {
			WriteJson(w, http.StatusForbidden, Error{Error: err.Error()})
			return
		}
		if databaseUser.Username != claimsUser.Username ||
			databaseUser.Password.HashedPassword != claimsUser.Password.HashedPassword ||
			!databaseUser.Password.Active {
			WriteJson(w, http.StatusForbidden, fmt.Errorf("token not verified"))
			return
		}
		handlerfunc(w, r)
	}
}

func validateJwt(tokenString string) (*jwt.Token, error) {
	mySigninKey := os.Getenv("JWTKEY")

	return jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("Unexpected sighing method :%v", t.Header["alg"])
		}
		return []byte(mySigninKey), nil
	})
}
