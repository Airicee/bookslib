package main

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
    "os"

    _ "github.com/lib/pq"
    "golang.org/x/crypto/bcrypt" 
)

var db *sql.DB

type User struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func enableCORS(w *http.ResponseWriter) {
    (*w).Header().Set("Access-Control-Allow-Origin", "*")
    (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    (*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
    enableCORS(&w)
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }

    var u User
    if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    _, err = db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", u.Username, string(hashedPassword))
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
    enableCORS(&w)
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }

    var u User
    if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    var storedPassword string
    query := "SELECT password FROM users WHERE username = $1"
    err := db.QueryRow(query, u.Username).Scan(&storedPassword)

    if err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]string{"error": "Invalid credentials"})
        return
    }

    err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(u.Password))
    if err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]string{"error": "Invalid credentials"})
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"token": "dummy-jwt-token-for-" + u.Username})
}

func main() {
    var err error
    db, err = sql.Open("postgres", os.Getenv("DB_DSN"))
    if err != nil {
        log.Fatal(err)
    }

    _, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, username TEXT UNIQUE, password TEXT)")
    if err != nil {
        log.Fatal(err)
    }

    adminPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
    db.Exec("INSERT INTO users (username, password) VALUES ('admin', $1) ON CONFLICT DO NOTHING", string(adminPassword))

    http.HandleFunc("/login", loginHandler)
    http.HandleFunc("/register", registerHandler)

    log.Println("Auth service running on port 8081")
    if err := http.ListenAndServe(":8081", nil); err != nil {
        log.Fatal(err)
    }
}