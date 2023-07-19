package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetRealtimeWeather(t *testing.T) {
	handler := &Handler{
		Location: "30643",
	}

	req, err := http.NewRequest("GET", "/realtime-weather", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handlerFunc := http.HandlerFunc(handler.getRealtimeWeather)

	handlerFunc.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := `location parameter is required`
	if strings.TrimSpace(rr.Body.String()) != strings.TrimSpace(expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}
