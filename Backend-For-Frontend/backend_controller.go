package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func main() {

	fmt.Println("Starting the application...")
	handler := &Handler{
		Location: "30643",
		//databse will go here
	}

	http.HandleFunc("/weather-forecast", func(w http.ResponseWriter, r *http.Request) {
		handler.getForecast(w, r)
	})

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))

}

func (h *Handler) getForecast(w http.ResponseWriter, r *http.Request) {
	// Get location parameter from query string
	locationName := r.URL.Query().Get("location")

	fmt.Println(locationName)

	if locationName == "" {
		http.Error(w, "Missing 'location' parameter", http.StatusBadRequest)
		return
	}

	// Create the forecast URL
	forecastURL := fmt.Sprintf("http://127.0.0.1:8081/get-forecast?location=%s", url.QueryEscape(locationName))

	// Make the GET request
	resp, err := http.Get(forecastURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If response is 404, do something else
	if resp.StatusCode == http.StatusNotFound {
		// Do something else here
		fmt.Println("404")
		return
	} else if resp.StatusCode != http.StatusOK {
		http.Error(w, "Unexpected status code", http.StatusInternalServerError)
		return
	}

	// Decode JSON from response body
	var forecastData interface{}
	if err := json.NewDecoder(resp.Body).Decode(&forecastData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp.Body.Close()

	// Return JSON to user
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(forecastData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type Handler struct {
	//location variable from HTTP request
	Location string `json:"location"`
}

type Location struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Region  string  `json:"region"`
	Country string  `json:"country"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Url     string  `json:"url"`
}
