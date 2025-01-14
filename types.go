package main

type Auth struct {
	Username string   `json:"username"`
	Id       int64    `json:"id"`
	Password Password `json:"password"`
}

type Password struct {
	HashedPassword  string `json:"-_"`
	Active          bool   `json:"active"`
}

func CreateAuth(username string, password Password) *Auth {
	return &Auth{
		Username: username,
		Password: password,
	}
}
