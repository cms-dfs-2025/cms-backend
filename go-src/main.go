package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
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

func ConnectDB() (*sql.DB, error) {
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

	if err != nil {
		log.Print(err)
		return db, err
	}

	if err = db.Ping(); err != nil {
		log.Print(err)
		return db, err
	}

	return db, nil
}

func CreateServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		log.Print("Processing request on /users")

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = io.WriteString(w, "List of users!!")
	})

	port := os.Getenv("SERVER_PORT")
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	log.Printf("Server created on port %s\n", port)

	return server
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// database connection
	db, err := ConnectDB()
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	// http server creation
	server := CreateServer()
	log.Fatal(server.ListenAndServe())
}
