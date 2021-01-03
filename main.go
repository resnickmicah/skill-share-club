package main

import (
    "net/http"
    "log"
    // "github.com/gorilla/mux"
    "github.com/gomodule/redigo/redis"
)

// Store the redis connection as a package level variable
var cache redis.Conn

func main() {
    initCache()
    // r := mux.NewRouter()
    // "Signin" and "Welcome" are the handlers that we will implement
    http.HandleFunc("/", Root)
    http.HandleFunc("/signin", Signin)
    http.HandleFunc("/welcome", Welcome)
    // start the server on port 8000
    log.Fatal(http.ListenAndServe(":8000", nil))
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
