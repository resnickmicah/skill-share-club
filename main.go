package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gomodule/redigo/redis"
	// _ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
)

// Store the redis connection as a package level variable
var cache redis.Conn
var db *sql.DB

const port = ":8000"

func main() {
	initCache()
	// r := mux.NewRouter()
	// "Signin" and "Welcome" are the handlers that we will implement
	http.HandleFunc("/", Root)
	http.HandleFunc("/signin", Signin)
	http.HandleFunc("/signup", Signup)
	http.HandleFunc("/welcome", Welcome)
	http.HandleFunc("/refresh", Refresh)

	initDB()
	defer db.Close()
	// start the server on port 8000
	fmt.Printf("It's alive!!! Check localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func initCache() {
	// Initialize the redis connection to a redis instance running on your local machine
	conn, err := redis.DialURL("redis://localhost")
	if err != nil {
		panic(err)
	}
	// Assign the connection to the package level `cache` variable
	cache = conn
}

func initDB() {
	var err error
	// Connect to the postgres db
	//you might have to change the connection string to add your database credentials
	db, err = sql.Open("postgres", "dbname=skill-share-club sslmode=disable")
	if err != nil {
		panic(err)
	}
	println("Successfully connected to postgres!")
}
