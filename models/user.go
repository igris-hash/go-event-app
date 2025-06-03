package models

import (
	"errors"

	"github.com/igris-hash/go-event-app/db"
	"github.com/igris-hash/go-event-app/utils"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (u *User) Save() error {
	query := `
	INSERT INTO users(email, password) VALUES(?, ?)
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	hashPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}
	result, err := stmt.Exec(u.Email, hashPassword)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	u.ID = id
	return err

}

func (u *User) ValidateCredentials() error {
	query := `
	SELECT id, password FROM users WHERE email = ?
	`
	row := db.DB.QueryRow(query, u.Email)

	var retrievedPassword string
	err := row.Scan(&u.ID, &retrievedPassword)
	if err != nil {
		return errors.New("Credentials invalid")
	}

	isPasswordValid := utils.ComparePassword(retrievedPassword, u.Password)
	if !isPasswordValid {
		return errors.New("Credentials invalid")
	}
	return nil
}
