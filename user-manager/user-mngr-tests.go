package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSignup(t *testing.T) {
	// Create a new HTTP request with a JSON body containing the credentials
	reqBody := `{"username": "testuser", "password": "testpassword"}`
	req, err := http.NewRequest("POST", "/signup", strings.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a new HTTP recorder to capture the response
	rr := httptest.NewRecorder()

	// Call the Signup function with the request and recorder
	Signup(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body is what we expect
	expected := ""
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
