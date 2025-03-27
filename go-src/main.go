package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func LogUsers(db *sql.DB) error {
	res, err := db.Query("SELECT handle, is_admin, auth_bits FROM users")
	defer res.Close()

	if err != nil {
		return err
	}

	for res.Next() {
		var handle string
		var is_admin bool
		var auth_bits int

		if err = res.Scan(&handle, &is_admin, &auth_bits); err != nil {
			return err
		}

		log.Printf("handle=%s, is_admin=%t, auth_bits=%d",
			handle, is_admin, auth_bits)
	}

	if err = res.Err(); err != nil {
		return err
	}

	return nil
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	host := "localhost"
	port := os.Getenv("POSTGRES_CONNECT_PORT")
	user := os.Getenv("POSTGRES_USER")
	dbname := os.Getenv("POSTGRES_DBNAME")
	password := os.Getenv("POSTGRES_PASSWORD")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	log.Print("Connecting to the db")

	db, err := sql.Open("postgres", psqlInfo)
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	err = LogUsers(db)
	if err != nil {
		log.Fatal(err)
	}
}
