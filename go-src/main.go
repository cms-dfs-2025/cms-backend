package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"

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
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_CONNECT_PORT")
	user := os.Getenv("POSTGRES_USER")
	dbname := os.Getenv("POSTGRES_DBNAME")
	password := os.Getenv("POSTGRES_PASSWORD")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	log.Printf("Connecting to the db on port %s", port)

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

type signupBody struct {
	Handle   string `json:"handle"`
	Password string `json:"password"`
}

func (handler RequestHandler) HandleSignup(w http.ResponseWriter, r *http.Request) {
	log.Print("Processing /api/signup")

	if r.Method != http.MethodPost {
		log.Print("Method not allowed error")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var body signupBody
	err := decoder.Decode(&body)

	if err != nil {
		log.Print("Body json decoder error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = CreateUser(body.Handle, body.Password, false, handler.db)

	if err != nil {
		log.Print("User creation error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("User with handle %s created", body.Handle)
	w.WriteHeader(http.StatusOK)
	return
}

var authIncorrectFormatError = fmt.Errorf("Incorrect authorization format")
var ErrAuthFormat = errors.New("cms-server: incorrect auth format")

func decodeBase64(encoded []byte) ([]byte, error) {
	var decoded = make([]byte, base64.StdEncoding.DecodedLen(len(encoded)))
	n, err := base64.StdEncoding.Decode(decoded, encoded)

	if err != nil {
		return decoded, err
	}

	decoded = decoded[:n]
	return decoded, nil
}

// parses base64(base64(handle):password)
// returns handle, password, error
func parseBasicAuthorization(header_value string) (string, string, error) {
	if len(header_value) < 6 {
		return "", "", ErrAuthFormat
	}

	basic := header_value[:6]

	if basic != "Basic " {
		return "", "", ErrAuthFormat
	}

	passBytes, err := decodeBase64([]byte(header_value[6:]))

	if err != nil {
		return "", "", err
	}

	colon_idx := slices.Index(passBytes, ':')

	if colon_idx == -1 {
		return "", "", ErrAuthFormat
	}

	handleBytes, err := decodeBase64(passBytes[:colon_idx])

	if err != nil {
		return "", "", err
	}

	passwordBytes := passBytes[(colon_idx + 1):]

	return string(handleBytes), string(passwordBytes), nil
}

func (handler RequestHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	log.Print("Processing /api/login")

	if r.Method != http.MethodPost {
		log.Print("Method not allowed error")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	auth := r.Header.Get("Authorization")
	handle, password, err := parseBasicAuthorization(auth)

	if err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Basic authorization error")
		return
	}

	err = AuthorizeUser(handle, password, handler.db)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func CreateServer(db *sql.DB) *http.Server {
	requestHandler := RequestHandler{db: db}

	mux := http.NewServeMux()

	// --- endpoints creation ---
	mux.HandleFunc("/api/signup", requestHandler.HandleSignup)
	mux.HandleFunc("/api/login", requestHandler.HandleLogin)

	// --------------------------

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
