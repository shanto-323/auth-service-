package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type api struct {
	IPAddr string
	Store  Storage
}

func WriteError(w http.ResponseWriter, status int, m any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(m)
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
	router.HandleFunc("/create", handlerFunc(a.VerifyAccount))

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
			WriteError(w, http.StatusOK, err)
		}
	}
}

func (a *api) CreateAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	if len(username) < 8 && len(password) < 8 {
		return WriteError(w, http.StatusNotAcceptable, "username & password not allowed")
	}

	accExist, err := a.Store.VerifyAccount(username, password)
	if err != nil {
		return WriteError(w, http.StatusBadGateway, err)
	}
	// If account already exists in database
	if accExist == true {
		return WriteError(w, http.StatusMisdirectedRequest, "accout already exists")
	}

	hashPassword, err := HashPassword(password)
	if err != nil {
		return WriteError(w, http.StatusBadRequest, err)
	}

	authUser := CreateAuth(username, hashPassword)
	err = a.Store.CreateAccount(authUser)
	if err != nil {
		return WriteError(w, http.StatusBadRequest, err)
	}

	return WriteError(w, http.StatusOK, authUser)
}

// ie login into account
func (a *api) VerifyAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	if len(username) < 8 || len(password) < 8 {
		return WriteError(w, http.StatusNotAcceptable, "username & password not allowed")
	}

	accExist, err := a.Store.VerifyAccount(username, password)
	switch accExist {
	case false:
		if err != nil {
			return WriteError(w, http.StatusInternalServerError, "Database Error")
		} else {
			return WriteError(w, http.StatusBadRequest, "username not found")
		}
	case true:
		if err != nil {
			return WriteError(w, http.StatusBadGateway, "username found but password not matched")
		} else {
			return WriteError(w, http.StatusOK, "password matched")
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
