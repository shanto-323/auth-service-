package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Storage interface {
	CreateAccount(memberID *Auth) (*Auth, error)
	GetAccountById(id int) (*Auth, error)
	VerifyAccount(userName string, password string) (*Auth, bool, error)
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

	return &Db{db: db}, nil
}

func (db *Db) init() error {
	quary := `  
    CREATE TABLE IF NOT EXISTS auth (
      id SERIAL PRIMARY KEY,
      username        VARCHAR(255)  NOT NULL UNIQUE,
      password        VARCHAR(255)  NOT NULL,
      active          BOOLEAN
    )
  `
	_, err := db.db.Exec(quary)
	if err != nil {
		return err
	}
	return nil
}

func (db *Db) CreateAccount(memberID *Auth) (*Auth, error) {
	quary := `  
    INSERT INTO auth (username ,password ,active)
    VALUES(?,?,?) 
  `
	rows, err := db.db.Exec(
		quary,
		memberID.Username,
		memberID.Password.HashedPassword,
		memberID.Password.Active,
	)
	if err != nil {
		return nil, err
	}

	lastInserId, err := rows.LastInsertId()
	if err != nil {
		return nil, err
	}
	return db.GetAccountById(int(lastInserId))
}

func (db *Db) VerifyAccount(username string, password string) (*Auth, bool, error) {
	user := &Auth{}
	rows, err := db.db.Query("SELECT * FROM auth WHERE username=?", username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		if user, err = getAccount(rows); err != nil {
			return nil, true, err
		}
	}

	if err = checkPassword(password, user.Password.HashedPassword); err != nil {
		return nil, true, fmt.Errorf("password not match: %w", err)
	}
	return user, true, nil
}

func (db *Db) UpdateAccount(memberID *Auth) error {
	return nil
}

func (db *Db) DeleteAccount(id int) error {
	return nil
}

func (db *Db) GetAccountById(id int) (*Auth, error) {
	rows, err := db.db.Query("SELECT * FROM auth WHERE id=?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		return getAccount(rows)
	}
	return nil, fmt.Errorf("account %v not found", id)
}

func getAccount(rows *sql.Rows) (*Auth, error) {
	userAccount := &Auth{}
	err := rows.Scan(
		&userAccount.Id,
		&userAccount.Username,
		&userAccount.Password.HashedPassword,
		&userAccount.Password.Active)
	if err != nil {
		return nil, err
	}
	return userAccount, nil
}
