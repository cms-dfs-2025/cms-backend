package main

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"database/sql"
)

const (
	hashCost int = bcrypt.DefaultCost
)

func truncatePassword(password string) []byte {
	passwordBytes := []byte(password)
	if len(passwordBytes) > 72 {
		passwordBytes = passwordBytes[:72]
	}

	return passwordBytes
}

func createPasswordHash(password string) (string, error) {
	passwordBytes := truncatePassword(password)
	hashBytes, err := bcrypt.GenerateFromPassword(passwordBytes, hashCost)

	if err != nil {
		return "error", err
	}

	return string(hashBytes), nil
}

func compareHashPassword(hash string, password string) error {
	passwordBytes := truncatePassword(password)
	hashBytes := []byte(hash)
	return bcrypt.CompareHashAndPassword(hashBytes, passwordBytes)
}

var ErrUserExists = errors.New("cms: user exists")

func CreateUser(handle string, password string, is_admin bool, db *sql.DB) error {
	user_check := db.
		QueryRow("SELECT FROM users WHERE handle = $1;", handle)

	err := user_check.Scan()

	if err != sql.ErrNoRows {
		if err == nil {
			return ErrUserExists
		} else {
			return err
		}
	}

	hash, err := createPasswordHash(password)

	if err != nil {
		return err
	}

	result, err := db.
		Exec("INSERT INTO users (handle, is_admin, auth_bits) "+
			"VALUES ($1, $2, $3);", handle, is_admin, hash)

	if err != nil {
		return err
	}

	n, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if n != 1 {
		return errors.New("cms: insert user failed")
	}

	return nil
}

func ChangeUserPassword(handle string, new_password string, db *sql.DB) error {
	hash, err := createPasswordHash(new_password)

	if err != nil {
		return err
	}

	result, err := db.
		Exec("UPDATE users SET auth_bits = $1 WHERE handle = $2;", hash, handle)

	if err != nil {
		return err
	}

	n, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if n != 1 {
		return errors.New("cms: update user password failed")
	}

	return nil
}

type UserRow struct {
	Id       int
	Handle   string
	IsAdmin  bool
	AuthBits string
}

func GetUser(handle string, db *sql.DB) (UserRow, error) {
	result := db.QueryRow(
		"SELECT id, handle, is_admin, auth_bits FROM users WHERE handle = $1;", handle)

	var user UserRow

	err := result.Scan(&user.Id, &user.Handle, &user.IsAdmin, &user.AuthBits)

	if err != nil {
		user.Id = -1
		user.Handle = ""
		user.IsAdmin = false
		user.AuthBits = "+"
	}

	return user, err
}

func AuthorizeUser(handle string, password string, db *sql.DB) (UserRow, error) {
	user, err := GetUser(handle, db)

	if err != nil {
		return user, err
	}

	err = compareHashPassword(user.AuthBits, password)
	return user, err
}
