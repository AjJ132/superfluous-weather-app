package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"golang.org/x/crypto/bcrypt"
)

type mockHasher struct{}

func (mh mockHasher) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return []byte("hashedPassword"), nil
}

// test signup
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
	mock.ExpectExec("INSERT INTO users").WithArgs(sql.Named("p1", "testuser"), sql.Named("p2", "hashedPassword")).WillReturnResult(sqlmock.NewResult(1, 1))

	// Create a Signup request
	creds := &Credentials{
		Password: "password",
		Username: "testuser",
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

// test login// test successful login
func TestSigninSuccess(t *testing.T) {
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
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	mock.ExpectQuery("SELECT sPassword FROM users WHERE sUsername=@username").WithArgs(sql.Named("username", "testuser")).WillReturnRows(sqlmock.NewRows([]string{"sPassword"}).AddRow(string(hashedPassword)))

	// Create a Signin request
	creds := &Credentials{
		Password: "password",
		Username: "testuser",
	}
	credsJson, _ := json.Marshal(creds)
	req, err := http.NewRequest("POST", "/signin", bytes.NewBuffer(credsJson))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the Signin handler
	handler.Signin(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := `{"message":"success"}`
	actual := strings.TrimSpace(rr.Body.String())
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}

	// Make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// test incorrect password
func TestSigninIncorrectPassword(t *testing.T) {
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
	mock.ExpectQuery("SELECT sPassword FROM users WHERE sUsername=@username").WithArgs(sql.Named("username", "testuser")).WillReturnRows(sqlmock.NewRows([]string{"sPassword"}).AddRow("hashedPassword"))

	// Create a Signin request
	creds := &Credentials{
		Password: "wrongpassword",
		Username: "testuser",
	}
	credsJson, _ := json.Marshal(creds)
	req, err := http.NewRequest("POST", "/signin", bytes.NewBuffer(credsJson))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the Signin handler
	handler.Signin(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	// Make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// test user not found
func TestSigninUserNotFound(t *testing.T) {
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
	mock.ExpectQuery("SELECT sPassword FROM users WHERE sUsername=@username").WithArgs(sql.Named("username", "testuser")).WillReturnError(sql.ErrNoRows)

	// Create a Signin request
	creds := &Credentials{
		Password: "password",
		Username: "testuser",
	}
	credsJson, _ := json.Marshal(creds)
	req, err := http.NewRequest("POST", "/signin", bytes.NewBuffer(credsJson))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the Signin handler
	handler.Signin(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	// Make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
