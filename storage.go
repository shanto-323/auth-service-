package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Storage interface {
	CreateAccount(memberID *Auth,refresh_token string) error
	VerifyAccount(userName string, password string) (bool, error)
	UpdateAccount(memberID *Auth) error
	DeleteAccount(id int) error
}

type Db struct {
	db *sql.DB
}

func CreateDb() (*Db, error) {
	dsn := "user1:12345678@tcp(localhost:3306)/mydb"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Db{
		db: db,
	}, nil
}

func (db *Db) init() error {
	quary := `  
    CREATE TABLE IF NOT EXISTS auth (
      id SERIAL PRIMARY KEY,
      username        VARCHAR(255)  NOT NULL UNIQUE,
      password        VARCHAR(255)  NOT NULL,
      refresh_token   TEXT          NOT NULL,
      active          BOOLEAN
    )
  `
	_, err := db.db.Exec(quary)
	if err != nil {
		return err
	}
	return nil
}

func (db *Db) CreateAccount(memberID *Auth, refresh_token string) error {
	quary := `  
    INSERT INTO auth (username , password, refresh_token,active)
    VALUES(?,?,?,?) 
  `
	_, err := db.db.Query(
		quary, memberID.Username,
		memberID.Password.HashedPassword,
		refresh_token,
		memberID.Password.Active,
	)
	if err != nil {
		return err
	}
	return nil
}

func (db *Db) VerifyAccount(username string, password string) (bool, error) {
	var auth Auth
	quary := `  
    SELECT id , username , password FROM auth WHERE username = ?
  `
	row := db.db.QueryRow(quary, username)
	err := row.Scan(
		&auth.Id,
		&auth.Username,
		&auth.Password.HashedPassword,
		&auth.Password.Active, //Chech if user is active (later)
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// Return false if data not exist
			return false, nil
		}
		// Database error
		return false, err
	}
	err = checkPassword(password, auth.Password.HashedPassword)
	if err != nil {
		// Username exists but password not match
		return true, err
	}
	// Return true if data exist and password matched
	return true, nil
}

func (db *Db) UpdateAccount(memberID *Auth) error {
	return nil
}

func (db *Db) DeleteAccount(id int) error {
	return nil
}
