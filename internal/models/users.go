package models

import (
	"database/sql"
	"errors"
	"time"
	"golang.org/x/crypto/bcrypt"

)

type User struct {
	ID           int
	Username     string
	PasswordHash []byte
	CreatedAt    time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Authenticate(username, password string) (int, error) {
	var id int
	var passwordHash []byte

	stmt := "SELECT id, password_hash FROM users WHERE username = ?"
	row := m.DB.QueryRow(stmt, username)

	err := row.Scan(&id, &passwordHash) 
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("invalid credentials")
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(passwordHash, []byte(password))
	if err != nil {
		return 0, errors.New("invalid credentials")
	}

	return id, nil 
}

func (m *UserModel) Insert(username string, passwordHash []byte) (int, error) {
	//writing the sql statement
	stmt := "INSERT INTO users (username, password_hash, created_at) VALUES (?, ?, UTC_TIMESTAMP())"

	//running the query on the db
	result, err := m.DB.Exec(stmt, username, passwordHash)
	if err != nil {
		return 0, err
	}

	//getting the users id from the db
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *UserModel) UserExists(username string) (bool, error) {
	stmt := `SELECT id FROM users WHERE username = ?`

	var count int
	err := m.DB.QueryRow(stmt, username).Scan(&count)
	if err != nil {
		return false, err
	}

	return true, nil 
}