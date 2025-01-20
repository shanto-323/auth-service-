package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

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
	router.HandleFunc("/login", handlerFunc(a.Login))
	router.HandleFunc("/account/{id}", withJwt(handlerFunc(a.GetAccountById), a.Store))

	fmt.Println("Api running...")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}

func handlerFunc(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJson(w, http.StatusBadRequest, Error{Error: err.Error()})
		}
	}
}

func (a *api) CreateAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return WriteJson(w, http.StatusMethodNotAllowed, Error{Error: "method not allowed"})
	}
	req := &Request{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return WriteJson(w, http.StatusBadRequest, req)
	}
	if err := checkLength(req.Username, req.Password); err != nil {
		return WriteJson(w, http.StatusNotAcceptable, Error{Error: err.Error()})
	}

	_, exist, err := a.Store.VerifyAccount(req.Username, req.Password)
	if err != nil && !exist {
		return WriteJson(w, http.StatusNotAcceptable, Error{Error: "account already exists"})
	}

	hashPassword, err := HashPassword(req.Password)
	if err != nil {
		return WriteJson(w, http.StatusBadRequest, Error{Error: err.Error()})
	}

	user := CreateAuth(req.Username, Password{HashedPassword: hashPassword, Active: true})
	//Save account in database
	account, err := a.Store.CreateAccount(user)
	if err != nil {
		return WriteJson(w, http.StatusBadRequest, Error{Error: err.Error()})
	}
	//Create jwt token
	jwtToken, err := createJWT(account)
	if err != nil {
		return WriteJson(w, http.StatusBadRequest, Error{Error: err.Error()})
	}
	return WriteJson(w, http.StatusOK, Response{
		Username: account.Username,
		Token:    jwtToken,
	})
}

func (a *api) Login(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return WriteJson(w, http.StatusBadRequest, Error{Error: "method not allowed"})
	}
	req := &Request{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return WriteJson(w, http.StatusBadRequest, err.Error())
	}
	if err := checkLength(req.Username, req.Password); err != nil {
		return WriteJson(w, http.StatusNotAcceptable, Error{Error: err.Error()})
	}
	account, exist, err := a.Store.VerifyAccount(req.Username, req.Password)
	if err != nil && !exist {
		return WriteJson(w, http.StatusNotAcceptable, Error{Error: "user not found"})
	}
	//Create jwt token
	jwtToken, err := createJWT(account)
	if err != nil {
		return WriteJson(w, http.StatusBadRequest, Error{Error: err.Error()})
	}
	return WriteJson(w, http.StatusOK, Response{
		Username: account.Username,
		Token:    jwtToken,
	})
}

func (a *api) GetAccountById(w http.ResponseWriter, r *http.Request) error {
	idString := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		return WriteJson(w, http.StatusBadRequest, err.Error())
	}
	user, err := a.Store.GetAccountById(id)
	if err != nil {
		return WriteJson(w, http.StatusBadRequest, err.Error())
	}
	return WriteJson(w, http.StatusOK, user)
}

func (a *api) UpdateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (a *api) DeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}
