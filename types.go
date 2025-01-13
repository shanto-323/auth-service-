package main

type Auth struct {
	Username       string `json:"string"`
	Id             int64  `json:"id"`
	HashedPassword string `json:"password"`
}

func CreateAuth(username string, hashedPassword string) *Auth {
	return &Auth{
		Username:       username,
		HashedPassword: hashedPassword,
	}
}
