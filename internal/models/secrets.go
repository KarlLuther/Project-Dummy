package models

import (
	"database/sql"
	"errors"
	"time"
)

type Secret struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Name        string    `json:"name"`
	SecretData  []byte    `json:"secret_data"`
	Created     time.Time `json:"created"`
	Expires     time.Time `json:"expires"`
}

type SecretModel struct {
	DB *sql.DB 
}

func (m *SecretModel) Insert(user_id int, name string, secret_data []byte, expires int)  (int, error){
	stmt := `INSERT INTO secrets(user_id, name, secret_data, created, expires)
	VALUES (?, ?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`


	//this command returns a sql.result object which has two methods: 
	//rowsAffected() and LasrInsertId()
	result, err := m.DB.Exec(stmt, user_id, name, secret_data, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	//id is of type int64 so we convert it into int before parsing forward
	return int(id), nil
}


func (m *SecretModel) Get(id int, userID int) (Secret, error) {
	stmt := `SELECT id, user_id,  name, secret_data, created, expires FROM secrets
	WHERE id = ? AND user_id = ? AND expires > UTC_TIMESTAMP()`

	row := m.DB.QueryRow(stmt,id,userID)

	var secret Secret

	err := row.Scan(&secret.ID,&secret.UserID, &secret.Name, &secret.SecretData, &secret.Created, &secret.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Secret{}, ErrNoRecord
		} else {
			return Secret{}, err
		}
	}

	return secret, nil
}