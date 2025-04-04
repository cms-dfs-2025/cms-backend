package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func GetUsers(db *sql.DB) ([]string, error) {
	res, err := db.Query("SELECT handle FROM users")
	defer res.Close()

	if err != nil {
		return nil, err
	}

	ret := []string{}
	for res.Next() {
		var handle string

		if err = res.Scan(&handle); err != nil {
			return nil, err
		}

		ret = append(ret, handle)
	}

	if err = res.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}

func ConnectDB() (*sql.DB, error) {
	host := "database"
	port := "5432"
	user := os.Getenv("POSTGRES_USER")
	dbname := os.Getenv("POSTGRES_DBNAME")
	password := os.Getenv("POSTGRES_PASSWORD")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	log.Println(psqlInfo)

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

type RequestHandler struct {
	db *sql.DB
}

func (handler RequestHandler) HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	log.Print("Processing /users")

	users, err := GetUsers(handler.db)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if err != nil {
		log.Print("Internal server error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	data := strings.Join(users, ",")

	_, _ = io.WriteString(w, data)
}

func CreateServer(db *sql.DB) *http.Server {
	requestHandler := RequestHandler{db: db}

	mux := http.NewServeMux()
	mux.HandleFunc("/users", requestHandler.HandleGetUsers)

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
	server := CreateServer(db)
	log.Fatal(server.ListenAndServe())
}
