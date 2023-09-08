package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func main() {

	fmt.Println("Starting the application...")
	handler := &Handler{
		Location: "30643",
		//databse will go here
	}

	http.HandleFunc("/api/weather-forecast", func(w http.ResponseWriter, r *http.Request) {
		handler.getForecast(w, r)
	})
	http.HandleFunc("/api/signin", func(w http.ResponseWriter, r *http.Request) {
		handler.signin(w, r)
	})
	http.HandleFunc("/api/signup", func(w http.ResponseWriter, r *http.Request) {
		handler.signup(w, r)
	})
	http.HandleFunc("/api/hello-world", func(w http.ResponseWriter, r *http.Request) {
		handler.helloWorld(w, r)
	})

	log.Fatal(http.ListenAndServe("0.0.0.0:8082", nil))

}

func (h *Handler) signin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Signin Received")
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
	resp, err := http.Post("http://login-service:8085/signin", "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//print response
	fmt.Println(resp)

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "failed to sign in", http.StatusUnauthorized)
		return
	}

	// Decode token from login system response
	var tokenResponse struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create JSON response with token
	response := struct {
		Status int    `json:"status"`
		Token  string `json:"token"`
	}{
		Status: http.StatusOK,
		Token:  tokenResponse.Token,
	}

	// Write JSON response to client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) signup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Signup Received")
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
	resp, err := http.Post("http://login-service:8085/signup", "application/json", bytes.NewBuffer(jsonBytes))
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

	// Decode token from login system response
	var tokenResponse struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create JSON response with token
	response := struct {
		Status int    `json:"status"`
		Token  string `json:"token"`
	}{
		Status: http.StatusOK,
		Token:  tokenResponse.Token,
	}

	// Write JSON response to client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Testing request return Hellow world as string
func (h *Handler) helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello World Received")
	// Set the appropriate headers for CORS
	w.Header().Set("Access-Control-Allow-Origin", "*") // replace '*' with a specific origin if needed
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	rand.Seed(time.Now().UnixNano())
	message := "Hello World: " + strconv.Itoa(rand.Intn(100))

	fmt.Println("Wrote Message")

	response := map[string]string{"message": message}

	w.Header().Set("Content-Type", "application/json")

	// Set the status code here, before writing the response
	w.WriteHeader(http.StatusOK)

	fmt.Println("Status OK")

	// Then encode the JSON response
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) getForecast(w http.ResponseWriter, r *http.Request) {
	
	startTime := time.Now()
	fromCache := false
	
	// Get location parameter from query string
	locationName := r.URL.Query().Get("location")

	if locationName == "" {
		http.Error(w, "Missing 'location' parameter", http.StatusBadRequest)
		return
	}

	// Create the forecast URL
	forecastURL := fmt.Sprintf("http://super-weather-cache-service:8091/get-forecast?location=%s", url.QueryEscape(locationName))

	// Make the GET request to cache
	resp, err := http.Get(forecastURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If response is 404, Fetch the weather
	if resp.StatusCode == http.StatusNotFound {
		forecastRequest := fmt.Sprintf("http://weather-mngr-service:8094/realtime-weather?location=%s", url.QueryEscape(locationName))
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

		fmt.Println("Data:", forecastData)

		response.Body.Close()

		fmt.Println("Successfully decoded weather data")

		// Prepare data to be sent to the save-forecast endpoint
		saveData, err := json.Marshal(forecastData)
		if err != nil {
			http.Error(w, "Error preparing data for save endpoint", http.StatusInternalServerError)
			return
		}

		// Send the forecast data to the save-forecast endpoint
		saveReq, err := http.NewRequest("POST", "http://super-weather-cache-service:8091/save-forecast", bytes.NewBuffer(saveData))
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

		endTime := time.Now()
		responseTime := endTime.Sub(startTime)

		// Create the map for meta data
		meta := map[string]interface{}{
			"responseTime": responseTime.Milliseconds(), // or .Seconds(), depending on what you want
			"fromCache":    fromCache, // Boolean value to indicate if data is from cache
		}

		// Directly add the 'meta' field to the existing forecastData map
		forecastData["meta"] = meta

		// Continue with sending the JSON response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(forecastData); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}


		return

	} else if resp.StatusCode != http.StatusOK {
		http.Error(w, "Unexpected status code", http.StatusInternalServerError)
		return
	}

	fmt.Println("Successfully got weather data from cache. Reutrning data to user")
	
	// Set fromCache to true
	fromCache = true
	endTime := time.Now()
	responseTime := endTime.Sub(startTime)

	// Prepare metadata and append it to forecastData
	meta := map[string]interface{}{
		"responseTime": responseTime.Milliseconds(), // or .Seconds(), based on what you need
		"fromCache":    fromCache,
	}

	// Decode JSON from response body
	var forecastData interface{}
	if err := json.NewDecoder(resp.Body).Decode(&forecastData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Assume forecastData is of type interface{}
	forecastDataMap, ok := forecastData.(map[string]interface{})
	if !ok {
		// Handle type assertion failure
		http.Error(w, "Type assertion failed", http.StatusInternalServerError)
		return
	}

	// Now you can set the 'meta' field
	forecastDataMap["meta"] = meta

	resp.Body.Close()

	// Return JSON to user
	// Send the final response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(forecastDataMap); err != nil {
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
