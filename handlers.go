package main

import (
    // "encoding/json"
    "fmt"
    "net/http"
    "time"

    "github.com/satori/go.uuid"
    "github.com/gorilla/schema"
)

var users = map[string]string{
    "user1": "Password1",
    "user2": "Password2",
}

var schemaDecoder = schema.NewDecoder()

// Create a struct that models the structure of a user, both in the request body, and in the DB
type Credentials struct {
    Password string `schema:"password"`
    Username string `schema:"username"`
}

func Root(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "./views/index.html")
}

func Signin(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w,"gorilla/schema couldn't parse the form: %s", err)
    }
    var creds Credentials
    // Get the JSON body and decode into credentials
    err = schemaDecoder.Decode(&creds, r.PostForm)
    if err != nil {
        // If the structure of the body is wrong, return an HTTP error
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w,"Map values do not match Credentials struct: %s", err)
        return
    }

    // Get the expected password from our in memory map
    expectedPassword, ok := users[creds.Username]

    // If a password exists for the given user
    // AND, if it is the same as the password we received, the we can move ahead
    // if NOT, then we return an "Unauthorized" status
    if !ok || expectedPassword != creds.Password {
        w.WriteHeader(http.StatusUnauthorized)
        fmt.Fprintf(w,"Your credentials are bad and you should feel bad: %s, %s", creds.Username, creds.Password)
        return
    }

    // Create a new random session token
    u, err := uuid.NewV4()
    if err != nil {
        fmt.Printf("Something went wrong: %s", err)
        return
    }
    sessionToken := u.String()
    // Set the token in the cache, along with the user whom it represents
    // The token has an expiry time of 120 seconds
    _, err = cache.Do("SETEX", sessionToken, "120", creds.Username)
    if err != nil {
        // If there is an error in setting the cache, return an internal server error
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    // Finally, we set the client cookie for "session_token" as the session token we just generated
    // we also set an expiry time of 120 seconds, the same as the cache
    http.SetCookie(w, &http.Cookie{
        Name:    "session_token",
        Value:   sessionToken,
        Expires: time.Now().Add(120 * time.Second),
    })
    fmt.Fprintf(w,"You did it! This app is so secure!")
}

func Welcome(w http.ResponseWriter, r *http.Request) {
    // We can obtain the session token from the requests cookies, which come with every request
    c, err := r.Cookie("session_token")
    if err != nil {
        if err == http.ErrNoCookie {
            // If the cookie is not set, return an unauthorized status
            w.WriteHeader(http.StatusUnauthorized)
            fmt.Fprintf(w, "Hey! Get outta here! You don't have a cookie!")
            return
        }
        // For any other type of error, return a bad request status
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, "Get outta here! Why? %s", err)
        return
    }
    sessionToken := c.Value

    // We then get the name of the user from our cache, where we set the session token
    response, err := cache.Do("GET", sessionToken)
    if err != nil {
        // If there is an error fetching from cache, return an internal server error status
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Uhhh... we fucked up. -- %s", err)
        return
    }
    if response == nil {
        // If the session token is not present in cache, return an unauthorized error
        w.WriteHeader(http.StatusUnauthorized)
        fmt.Fprintf(w, "I didn't find anything in the session cache... you should check there")
        return
    }
    // Finally, return the welcome message to the user
    w.Write([]byte(fmt.Sprintf("Welcome %s!", response)))
}

func Refresh(w http.ResponseWriter, r *http.Request) {
    c, err := r.Cookie("session_token")
    if err != nil {
        if err == http.ErrNoCookie {
            w.WriteHeader(http.StatusUnauthorized)

            return
        }
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    sessionToken := c.Value

    response, err := cache.Do("GET", sessionToken)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    if response == nil {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }
    // The code uptil this point is the same as the first part of the `Welcome` route

    // Now, create a new session token for the current user
    u, err := uuid.NewV4()
    if err != nil {
        fmt.Printf("Something went wrong: %s", err)
        return
    }
    newSessionToken := u.String()
    _, err = cache.Do("SETEX", newSessionToken, "120", fmt.Sprintf("%s",response))
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    // Delete the older session token
    _, err = cache.Do("DEL", sessionToken)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    
    // Set the new token as the users `session_token` cookie
    http.SetCookie(w, &http.Cookie{
        Name:    "session_token",
        Value:   newSessionToken,
        Expires: time.Now().Add(120 * time.Second),
    })
}
