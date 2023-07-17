package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/denisenkom/go-mssqldb"
)

// The "db" package level variable will hold the reference to our database instance
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

	log.Fatal(http.ListenAndServe("127.0.0.1:8000", nil))
}

func initDB() {
	fmt.Println("Connecting to database...")
	var err error
	// Connect to the MSSQL db using Windows Authentication
	// you might have to change the server name and database name
	connString := "server=localhost;database=Superfluous_Weather;integrated security=true;"
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		panic(err)
	}

	// check the connection
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected!")
}

// Hasher interface for hashing passwords
type Hasher interface {
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
}

// BcryptHasher is a concrete implementation of Hasher
type BcryptHasher struct{}

func (bh BcryptHasher) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}

// Handler contains the dependencies of the Signup function
type Handler struct {
	DB     *sql.DB
	Hasher Hasher
}

// Create a struct that models the structure of a user, both in the request body, and in the DB
type Credentials struct {
	xnEntryID int    `json:"xnEntryID", db:"xnEntryID"` // This field is ignored by go-sql-driver
	Password  string `json:"password", db:"sPassword"`
	Username  string `json:"username", db:"sUsername"`
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	// Parse and decode the request body into a new `Credentials` instance
	creds := &Credentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, err := h.Hasher.GenerateFromPassword([]byte(creds.Password), 8)

	fmt.Println("creds.Username: " + creds.Username)
	fmt.Println("creds.Password: " + creds.Password)
	fmt.Println("hashedPassword: " + string(hashedPassword))

	// Next, insert the username, along with the hashed password into the database
	if _, err = h.DB.Exec("INSERT INTO users (sUsername, sPassword) VALUES (@p1, @p2)", sql.Named("p1", creds.Username), sql.Named("p2", string(hashedPassword))); err != nil {

		fmt.Println("Error inserting into database: " + err.Error())
		// If there is any issue with inserting into the database, return a 500 error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// We reach this point if the credentials we correctly stored in the database, and the default status of 200 is sent back
}

func (h *Handler) Signin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Signin called")
	creds := &Credentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println("Querying database for username: " + creds.Username)
	result := h.DB.QueryRow("SELECT sPassword FROM users WHERE sUsername=@username", sql.Named("username", creds.Username))

	fmt.Println("Completed query")
	storedCreds := &Credentials{}
	err = result.Scan(&storedCreds.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println("storedCreds.Password: " + storedCreds.Password)

	// using the bcrypt hasher from Handler
	err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fmt.Println("Successfully logged in")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "success"})
}
