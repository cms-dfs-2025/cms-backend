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

func CreatePasswordHash(password string) (string, error) {
	passwordBytes := truncatePassword(password)
	hashBytes, err := bcrypt.GenerateFromPassword(passwordBytes, hashCost)

	if err != nil {
		return "error", err
	}

	return string(hashBytes), nil
}

func CompareHashPassword(hash string, password string) error {
	passwordBytes := truncatePassword(password)
	hashBytes := []byte(hash)
	return bcrypt.CompareHashAndPassword(hashBytes, passwordBytes)
}

func CreateUser(handle string, password string, is_admin bool, db *sql.DB) error {
	hash, err := CreatePasswordHash(password)

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
		return errors.New("cms-server: insert user failed")
	}

	return nil
}

type UserRow struct {
	handle    string
	is_admin  bool
	auth_bits string
}

func GetUser(handle string, db *sql.DB) (UserRow, error) {
	result := db.QueryRow(
		"SELECT handle, is_admin, auth_bits FROM users WHERE handle = $1", handle)

	var user UserRow

	err := result.Scan(&user.handle, &user.is_admin, &user.auth_bits)

	if err != nil {
		user.handle = ""
		user.is_admin = false
		user.auth_bits = "+"
	}

	return user, err
}

func AuthorizeUser(handle string, password string, db *sql.DB) error {
	user, err := GetUser(handle, db)

	if err != nil {
		return err
	}

	return CompareHashPassword(user.auth_bits, password)
}
