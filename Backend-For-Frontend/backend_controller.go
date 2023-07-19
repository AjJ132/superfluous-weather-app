package main

import (
	"bytes"
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

	http.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		handler.signin(w, r)
	})

	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		handler.signup(w, r)
	})

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))

}

func (h *Handler) getForecast(w http.ResponseWriter, r *http.Request) {
	// Get location parameter from query string
	locationName := r.URL.Query().Get("location")

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

	// If response is 404, Fetch the weather
	if resp.StatusCode == http.StatusNotFound {
		forecastRequest := fmt.Sprintf("http://127.0.0.1:8082/forecast?location=%s", url.QueryEscape(locationName))
		fmt.Println("Data not cached. Now calling weather api for data...")

		// Make the GET request
		response, err := http.Get(forecastRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println("Sent GET request to weather api")

		//Ensure status code is OK
		if response.StatusCode != http.StatusOK {
			http.Error(w, "Unexpected status code", http.StatusInternalServerError)
			return
		}

		fmt.Println("Successfully got weather data")

		// Decode JSON from response body
		var forecastData map[string]interface{}
		if err := json.NewDecoder(response.Body).Decode(&forecastData); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response.Body.Close()

		fmt.Println("Successfully decoded weather data")

		// Prepare data to be sent to the save-forecast endpoint
		saveData, err := json.Marshal(forecastData)
		if err != nil {
			http.Error(w, "Error preparing data for save endpoint", http.StatusInternalServerError)
			return
		}

		// Send the forecast data to the save-forecast endpoint
		saveReq, err := http.NewRequest("POST", "http://localhost:8081/save-forecast", bytes.NewBuffer(saveData))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		saveReq.Header.Set("Content-Type", "application/json")

		fmt.Println("Sending data to save endpoint...")

		saveResp, err := http.DefaultClient.Do(saveReq)
		if err != nil {
			http.Error(w, "Error sending data to save endpoint", http.StatusInternalServerError)
			return
		}
		if saveResp.StatusCode != http.StatusOK {
			http.Error(w, "Error from save endpoint", saveResp.StatusCode)
			return
		}

		fmt.Println("Successfully sent data to save endpoint")

		// Return JSON to user
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(forecastData); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	} else if resp.StatusCode != http.StatusOK {
		http.Error(w, "Unexpected status code", http.StatusInternalServerError)
		return
	}

	fmt.Println("Successfully got weather data from cache. Reutrning data to user")

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

func (h *Handler) signin(w http.ResponseWriter, r *http.Request) {
	// Check for json objects 'username' and 'password'
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert credentials to JSON
	jsonBytes, err := json.Marshal(creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send request to login system
	resp, err := http.Post("http://localhost:8083/signin", "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "failed to sign in", http.StatusUnauthorized)
		return
	}

	//Write return message
	w.Write([]byte("Successfully signed in"))

	// Return success
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) signup(w http.ResponseWriter, r *http.Request) {
	// Check for json objects 'username' and 'password'
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert credentials to JSON
	jsonBytes, err := json.Marshal(creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send request to login system
	resp, err := http.Post("http://localhost:8083/signup", "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "failed to sign in", http.StatusUnauthorized)
		return
	}

	//Write return message
	w.Write([]byte("Successfully signed up"))

	// Return success
	w.WriteHeader(http.StatusOK)
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
