package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
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

	log.Fatal(http.ListenAndServe("0.0.0.0:8086", nil))
}

func initDB() {
	fmt.Println("Attempting to connect to database...")
	var err error

	connString := "postgres://dbsa:Admin123@login-database-service:5432/UserDatabase?sslmode=disable"
	db, err = sql.Open("postgres", connString)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Database connection successful!")
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
	XnEntryID int
	SPassword string `json:"password"`
	SUsername string `json:"username"`
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Signup request received")
	//decode request body into struct
	creds := &Credentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//hash password
	hashedPassword, err := h.Hasher.GenerateFromPassword([]byte(creds.SPassword), 10)

	//Check for errors from hashed password
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//insert user into database
	if _, err = h.DB.Exec(`INSERT INTO Users (sUsername, sPassword) VALUES ($1, $2)`, creds.SUsername, string(hashedPassword)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//generate JWT token
	token, err := generateToken(creds.SUsername)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

	//return token
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *Handler) Signin(w http.ResponseWriter, r *http.Request) {
	//decode request body into struct
	creds := &Credentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//query database for user
	result := h.DB.QueryRow(`SELECT sPassword FROM Users WHERE sUsername=$1`, creds.SUsername)
	storedCreds := &Credentials{}
	err = result.Scan(&storedCreds.SPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//compare passwords using bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(storedCreds.SPassword), []byte(creds.SPassword))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//generate JWT token
	token, err := generateToken(creds.SUsername)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

	//return token
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"token": token})
}

//Generate JWT Token
func generateToken(username string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": username,
        "exp":      time.Now().Add(time.Hour * 24).Unix(),
    })

	//return token 
    return token.SignedString([]byte("super-weather-secret-key"))
}
