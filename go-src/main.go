package main

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

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

const SERVER_CONTEXT = "ServerContext"

type ServerContext struct {
	db *sql.DB
}

type signupBody struct {
	Handle   string `json:"handle" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (handler ServerContext) HandleSignup(c *gin.Context) {
	var body signupBody
	err := c.ShouldBindJSON(&body)

	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = CreateUser(body.Handle, body.Password, false, handler.db)

	if err != nil {
		if err == ErrUserExists {
			c.JSON(http.StatusInternalServerError,
				gin.H{"message": "User already exists"})
		} else {
			c.JSON(http.StatusInternalServerError,
				gin.H{"message": "Internal error"})
		}
		return
	}

	c.Status(http.StatusOK)
	return
}

var ErrAuthFormat = errors.New("cms: incorrect auth format")

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

func (handler ServerContext) LoginMiddleware(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	handle, password, err := parseBasicAuthorization(auth)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			gin.H{"message": "Basic auth error"})
		return
	}

	user, err := AuthorizeUser(handle, password, handler.db)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			gin.H{"message": "Unauthorized"})
		return
	}

	c.Set("user", user)
}

func (handler ServerContext) HandleLogin(c *gin.Context) {
	c.Status(http.StatusOK)
	return
}

type changePwBody struct {
	Password string `json:"password" binding:"required"`
}

func (handler ServerContext) HandleChangePw(c *gin.Context) {
	user := c.MustGet("user").(UserRow)

	var body changePwBody
	err := c.ShouldBindJSON(&body)

	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = ChangeUserPassword(user.handle, body.Password, handler.db)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
	return
}

func CreateServer(db *sql.DB) *gin.Engine {
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	serverContext := &ServerContext{db: db}

	front_files := os.Getenv("FRONT_FILES")

	// --- endpoints creation ---
	router.Use(static.Serve("/", static.LocalFile(front_files, false)))

	router.POST("/api/auth/signup", serverContext.HandleSignup)
	router.POST("/api/auth/login", serverContext.LoginMiddleware,
		serverContext.HandleLogin)
	router.POST("/api/auth/change_pw", serverContext.LoginMiddleware,
		serverContext.HandleChangePw)

	// --------------------------

	return router
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

	port := os.Getenv("SERVER_PORT")
	address := fmt.Sprintf(":%s", port)

	log.Printf("Creating a server on port %s\n", port)

	log.Fatal(server.Run(address))
}
