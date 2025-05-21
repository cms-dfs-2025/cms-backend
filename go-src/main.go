package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

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

type ServerContext struct {
	db *sql.DB
}

func CreateServer(db *sql.DB) *gin.Engine {
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

	router.POST("/api/work/upload", serverContext.LoginMiddleware,
		serverContext.HandleWorkUpload)
	router.POST("/api/work/modidfy", serverContext.LoginMiddleware,
		serverContext.HandleWorkModify)
	router.POST("/api/work/delete", serverContext.LoginMiddleware,
		serverContext.HandleWorkDelete)
	router.GET("/api/work/get_all", serverContext.LoginMiddleware,
		serverContext.HandleWorkGetAll)
	router.GET("/api/work/get_body", serverContext.LoginMiddleware,
		serverContext.HandleWorkGetBody)

	router.GET("/api/posts/get_all", serverContext.HandlePostsGetAll)
	router.GET("/api/posts/get_body", serverContext.HandlePostsGetBody)

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
