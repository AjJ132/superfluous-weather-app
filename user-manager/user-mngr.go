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
	// "Signin" and "Signup" are handler that we will implement
	//http.HandleFunc("/signin", Signin)
	http.HandleFunc("/signup", Signup)
	// initialize our database connection
	fmt.Println("Starting the application...")
	initDB()
	fmt.Println("Connected to database")
	// start the server on port 8000
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

// Create a struct that models the structure of a user, both in the request body, and in the DB
type Credentials struct {
	Password string `json:"password", db:"sPassword"`
	Username string `json:"username", db:"sUsername"`
}

func Signup(w http.ResponseWriter, r *http.Request) {
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)

	fmt.Println("creds.Username: " + creds.Username)
	fmt.Println("creds.Password: " + creds.Password)
	fmt.Println("hashedPassword: " + string(hashedPassword))

	// Next, insert the username, along with the hashed password into the database
	if _, err = db.Exec("INSERT INTO users (sUsername, sPassword) VALUES (@p1, @p2)", sql.Named("p1", creds.Username), sql.Named("p2", string(hashedPassword))); err != nil {

		fmt.Println("Error inserting into database: " + err.Error())
		// If there is any issue with inserting into the database, return a 500 error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// We reach this point if the credentials we correctly stored in the database, and the default status of 200 is sent back
}
