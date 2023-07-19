package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"golang.org/x/crypto/bcrypt"
)

type mockHasher struct{}

func (m mockHasher) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return []byte("hashedPassword"), nil
}

func TestSignup(t *testing.T) {
	// Create a mock database
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	// Create a mock hasher
	hasher := mockHasher{}

	// Create the handler with the mock database and hasher
	handler := Handler{DB: mockDB, Hasher: hasher}

	// Define the expectations for the database interaction
	mock.ExpectExec("INSERT INTO Users").WithArgs("testuser", "hashedPassword").WillReturnResult(sqlmock.NewResult(1, 1))

	// Create a Signup request
	creds := &Credentials{
		SPassword: "password",
		SUsername: "testuser",
	}
	credsJson, _ := json.Marshal(creds)
	req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(credsJson))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the Signup handler
	handler.Signup(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSignin(t *testing.T) {
	// Create a mock database
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	fmt.Println("Connected to database")

	// Create a mock hasher
	hasher := mockHasher{}

	// Create the handler with the mock database and hasher
	handler := Handler{DB: mockDB, Hasher: hasher}

	fmt.Println("Connected!")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), 10)

	rows := sqlmock.NewRows([]string{"sPassword"}).AddRow(string(hashedPassword))

	mock.ExpectQuery("SELECT sPassword FROM Users WHERE sUsername=$1").WithArgs("testuser").WillReturnRows(rows)

	fmt.Println("Expectations defined")

	// Create a Signin request
	creds := &Credentials{
		SPassword: "password",
		SUsername: "testuser",
	}

	fmt.Println("Credentials created")

	credsJson, _ := json.Marshal(creds)
	req, err := http.NewRequest("POST", "/signin", bytes.NewBuffer(credsJson))
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Request created")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the Signin handler
	handler.Signin(rr, req)

	fmt.Println("Signin called")

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
