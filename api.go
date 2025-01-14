package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type api struct {
	IPAddr string
	Store  Storage
}

func MakeApi(ipAddr string, store Storage) *api {
	return &api{
		IPAddr: ipAddr,
		Store:  store,
	}
}

func (a *api) run() {
	router := chi.NewRouter()

	router.HandleFunc("/create", handlerFunc(a.CreateAccount))
	router.HandleFunc("/login", handlerFunc(a.VerifyAccount))

	fmt.Println("Api running...")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}

func handlerFunc(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			WriteJson(w, http.StatusBadRequest, Error{Error: err.Error()})
		}
	}
}

func (a *api) CreateAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return WriteJson(w, http.StatusMethodNotAllowed, Error{Error: "method not allowed"})
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	if len(username) < 8 || len(password) < 8 {
		return WriteJson(w, http.StatusNotAcceptable, Error{Error: "not enough length"})
	}

	accExist, err := a.Store.VerifyAccount(username, password) //Create another function for just verify the user (later)
	if err != nil {
		return WriteJson(w, http.StatusBadRequest, Error{Error: err.Error()})
	}
	// If account already exists in database
	if accExist == true {
		return WriteJson(w, http.StatusNotAcceptable, Error{Error: "account already exists"})
	}

	//Create hash password
	hashPassword, err := HashPassword(password)
	if err != nil {
		return WriteJson(w, http.StatusBadRequest, Error{Error: err.Error()})
	}

	//Create jwt token
	jwtToken, err := createJWT(username, hashPassword)
	if err != nil {
		return WriteJson(w, http.StatusBadRequest, Error{Error: err.Error()})
	}

	authUser := CreateAuth(username, Password{HashedPassword: hashPassword, Active: true})

	// Create refresh token
	refreshToken, err := createRefreshToken(authUser)
	if err != nil {
		return WriteJson(w, http.StatusBadRequest, Error{Error: err.Error()})
	}

	//Save account in database
	err = a.Store.CreateAccount(authUser, refreshToken)
	if err != nil {
		return WriteJson(w, http.StatusBadRequest, Error{Error: err.Error()})
	}

	return WriteJson(w, http.StatusOK, map[string]string{
		"jwt_token":     jwtToken,
		"refresh_token": refreshToken,
	})
}

// ie login into account
func (a *api) VerifyAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return WriteJson(w, http.StatusBadRequest, Error{Error: "Method Not Allowed"})
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	if len(username) < 8 || len(password) < 8 {
		return WriteJson(w, http.StatusBadRequest, Error{Error: "not enough lenght"})
	}

	accExist, err := a.Store.VerifyAccount(username, password)
	switch accExist {
	case false:
		if err != nil {
			return WriteJson(w, http.StatusInternalServerError, Error{Error: "Database Error"})
		} else {
			return WriteJson(w, http.StatusBadRequest, Error{Error: "User Not Exists"})
		}
	case true:
		if err != nil {
			return WriteJson(w, http.StatusBadRequest, Error{Error: "Wrong Password"})
		}
	}
	return nil
}

func (a *api) UpdateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (a *api) DeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}
