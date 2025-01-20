package main

type Auth struct {
	Id       int64    `json:"id"`
	Username string   `json:"username"`
	Password Password `json:"password"`
}

type Password struct {
	HashedPassword string `json:"-_"`
	Active         bool   `json:"active"`
}

func CreateAuth(username string, password Password) *Auth {
	return &Auth{
		Username: username,
		Password: password,
	}
}

type Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}
