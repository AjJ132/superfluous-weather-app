package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func main() {
	fmt.Println("Starting the application...")
	initDB()
	fmt.Println("Connected to database")

	handler := &Handler{
		DB:     db,
		Hasher: BcryptHasher{},
	}

	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		handler.Signup(w, r)
	})

	http.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		handler.Signin(w, r)
	})

	log.Fatal(http.ListenAndServe(":8083", nil))
}

func initDB() {
	fmt.Println("Connecting to database...")
	var err error

	connString := "postgres://dbsa:Admin123@db:5432/UserDatabase?sslmode=disable"
	db, err = sql.Open("postgres", connString)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected!")
}

type Hasher interface {
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
}

type BcryptHasher struct{}

func (bh BcryptHasher) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}

type Handler struct {
	DB     *sql.DB
	Hasher Hasher
}

type Credentials struct {
	xnEntryID int
	sPassword string `json:"password"`
	sUsername string `json:"username"`
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Signup request received")
	creds := &Credentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println("Username: " + creds.sUsername)
	fmt.Println("Password: " + creds.sPassword)
	hashedPassword, err := h.Hasher.GenerateFromPassword([]byte(creds.sPassword), 8)

	if _, err = h.DB.Exec(`INSERT INTO Users (sUsername, sPassword) VALUES ($1, $2)`, creds.sUsername, string(hashedPassword)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "success"})
}

func (h *Handler) Signin(w http.ResponseWriter, r *http.Request) {
	creds := &Credentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result := h.DB.QueryRow(`SELECT sPassword FROM Users WHERE sUsername=$1`, creds.sUsername)
	storedCreds := &Credentials{}
	err = result.Scan(&storedCreds.sPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedCreds.sPassword), []byte(creds.sPassword))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "success"})
}
