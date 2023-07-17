package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

//correct format for calling this script is below
//http://localhost:8080/realtime-weather?location=30144

func main() {

	fmt.Println("Starting the application...")
	handler := &Handler{
		Location: "30643",
		//databse will go here
	}

	http.HandleFunc("/realtime-weather", func(w http.ResponseWriter, r *http.Request) {
		handler.getRealtimeWeather(w, r)
	})

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}

type Handler struct {
	//location variable from HTTP request
	Location string `json:"location"`
}

func (h *Handler) getRealtimeWeather(w http.ResponseWriter, r *http.Request) {
	location := r.URL.Query().Get("location")
	if location == "" {
		http.Error(w, "location parameter is required", http.StatusBadRequest)
		return
	}
	url := fmt.Sprintf("https://weatherapi-com.p.rapidapi.com/current.json?q=%s", location)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", "2fdb8b6f67mshb1178afe669a190p1a700djsnc0d0e7bce537")
	req.Header.Add("X-RapidAPI-Host", "weatherapi-com.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var weather RealtimeWeather
	_ = json.Unmarshal(body, &weather)

	fmt.Println("Successfully got weather")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(weather)

}

type RealtimeWeather struct {
	Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Lat            float64 `json:"lat"`
		Lon            float64 `json:"lon"`
		TzID           string  `json:"tz_id"`
		LocaltimeEpoch int     `json:"localtime_epoch"`
		Localtime      string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		LastUpdatedEpoch int     `json:"last_updated_epoch"`
		LastUpdated      string  `json:"last_updated"`
		TempC            float64 `json:"temp_c"`
		TempF            float64 `json:"temp_f"`
		IsDay            int     `json:"is_day"`
		Condition        struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
			Code int    `json:"code"`
		} `json:"condition"`
		WindMph    float64 `json:"wind_mph"`
		WindKph    float64 `json:"wind_kph"`
		WindDegree int     `json:"wind_degree"`
		WindDir    string  `json:"wind_dir"`
		PressureMb float64 `json:"pressure_mb"`
		PressureIn float64 `json:"pressure_in"`
		PrecipMm   float64 `json:"precip_mm"`
		PrecipIn   float64 `json:"precip_in"`
		Humidity   int     `json:"humidity"`
		Cloud      int     `json:"cloud"`
		FeelslikeC float64 `json:"feelslike_c"`
		FeelslikeF float64 `json:"feelslike_f"`
		VisKm      float64 `json:"vis_km"`
		VisMiles   float64 `json:"vis_miles"`
		UV         int     `json:"uv"`
		GustMph    float64 `json:"gust_mph"`
		GustKph    float64 `json:"gust_kph"`
	} `json:"current"`
}
