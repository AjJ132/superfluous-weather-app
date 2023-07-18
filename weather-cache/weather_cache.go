package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/denisenkom/go-mssqldb"
)

//will hold refrence to database instance
var db *sql.DB

func main() {
	fmt.Println("Starting the application...")
	initDB()
	fmt.Println("Connected to database")

	handler := &Handler{
		DB:     db,
	}

	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		handler.getCurrentWeather(w, r)
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
		fmt.Println("Cache Serice Database Error")
		panic(err)
	}

	// check the connection
	err = db.Ping()
	if err != nil {
		fmt.Println("Cache Serice Database Error")
		panic(err)
	}

	fmt.Println("Connected!")
}

//Get the current weather from the database
func (h *Handler) getCurrentWeather(w http.ResponseWriter, r *http.Request) {
	
	//first check to see if the database has the current weather and is not expired

	//if expired then call the API
	//End of API call
}

// Handler contains the dependencies of the Signup function
type Handler struct {
	DB     *sql.DB
}