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
//http://localhost:8080/forecast?location=London

func main() {

	fmt.Println("Starting the application...")
	handler := &Handler{
		Location: "30643",
		//databse will go here
	}

	http.HandleFunc("/realtime-weather", func(w http.ResponseWriter, r *http.Request) {
		handler.getRealtimeWeather(w, r)
	})

	// http.HandleFunc("/forecast", func(w http.ResponseWriter, r *http.Request) {
	// 	handler.getFutureWeather(w, r)
	// })

	log.Fatal(http.ListenAndServe(":8093", nil))
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

func (h *Handler) getFutureWeather(w http.ResponseWriter, r *http.Request) {
	location := r.URL.Query().Get("location")
	if location == "" {
		http.Error(w, "location parameter is required", http.StatusBadRequest)
		return
	}
	url := fmt.Sprintf("https://weatherapi-com.p.rapidapi.com/forecast.json?q=%s&days=3", location)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", "e6552f5551msh8af99331b1d77c7p16e6c2jsn105a4cbb1fbe")
	req.Header.Add("X-RapidAPI-Host", "weatherapi-com.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var forecast WeatherForecast
	_ = json.Unmarshal(body, &forecast)

	fmt.Println("Successfully got forecast")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(forecast)
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

type WeatherForecast struct {
	Location Location `json:"location"`
	Current  Current  `json:"current"`
	Forecast Forecast `json:"forecast"`
}

type Location struct {
	Name           string  `json:"name"`
	Region         string  `json:"region"`
	Country        string  `json:"country"`
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	TzID           string  `json:"tz_id"`
	LocaltimeEpoch int64   `json:"localtime_epoch"`
	Localtime      string  `json:"localtime"`
}

type Current struct {
	LastUpdatedEpoch int       `json:"last_updated_epoch"`
	LastUpdated      string    `json:"last_updated"`
	TempC            float64   `json:"temp_c"`
	TempF            float64   `json:"temp_f"`
	IsDay            int       `json:"is_day"`
	Condition        Condition `json:"condition"`
	WindMph          float64   `json:"wind_mph"`
	WindKph          float64   `json:"wind_kph"`
	WindDegree       int       `json:"wind_degree"`
	WindDir          string    `json:"wind_dir"`
	PressureMb       float64   `json:"pressure_mb"`
	PressureIn       float64   `json:"pressure_in"`
	PrecipMm         float64   `json:"precip_mm"`
	PrecipIn         float64   `json:"precip_in"`
	Humidity         int       `json:"humidity"`
	Cloud            int       `json:"cloud"`
	FeelslikeC       float64   `json:"feelslike_c"`
	FeelslikeF       float64   `json:"feelslike_f"`
	VisKm            float64   `json:"vis_km"`
	VisMiles         float64   `json:"vis_miles"`
	UV               int       `json:"uv"`
	GustMph          float64   `json:"gust_mph"`
	GustKph          float64   `json:"gust_kph"`
}

type Forecast struct {
	Forecastday []ForecastDay `json:"forecastday"`
}

type ForecastDay struct {
	Date      string `json:"date"`
	DateEpoch int64  `json:"date_epoch"`
	Day       Day    `json:"day"`
	Astro     Astro  `json:"astro"`
	Hour      []Hour `json:"hour"`
}

type Day struct {
	MaxtempC          float64   `json:"maxtemp_c"`
	MaxtempF          float64   `json:"maxtemp_f"`
	MintempC          float64   `json:"mintemp_c"`
	MintempF          float64   `json:"mintemp_f"`
	AvgtempC          float64   `json:"avgtemp_c"`
	AvgtempF          float64   `json:"avgtemp_f"`
	MaxwindMph        float64   `json:"maxwind_mph"`
	MaxwindKph        float64   `json:"maxwind_kph"`
	TotalprecipMm     float64   `json:"totalprecip_mm"`
	TotalprecipIn     float64   `json:"totalprecip_in"`
	TotalsnowCm       float64   `json:"totalsnow_cm"`
	AvgvisKm          float64   `json:"avgvis_km"`
	AvgvisMiles       float64   `json:"avgvis_miles"`
	Avghumidity       float64   `json:"avghumidity"`
	DailyWillItRain   int       `json:"daily_will_it_rain"`
	DailyChanceOfRain string    `json:"daily_chance_of_rain"`
	DailyWillItSnow   int       `json:"daily_will_it_snow"`
	DailyChanceOfSnow string    `json:"daily_chance_of_snow"`
	Condition         Condition `json:"condition"`
	UV                float64   `json:"uv"`
}

type Astro struct {
	Sunrise          string `json:"sunrise"`
	Sunset           string `json:"sunset"`
	Moonrise         string `json:"moonrise"`
	Moonset          string `json:"moonset"`
	MoonPhase        string `json:"moon_phase"`
	MoonIllumination string `json:"moon_illumination"`
	IsMoonUp         int    `json:"is_moon_up"`
	IsSunUp          int    `json:"is_sun_up"`
}

type Hour struct {
	TimeEpoch    int       `json:"time_epoch"`
	Time         string    `json:"time"`
	TempC        float64   `json:"temp_c"`
	TempF        float64   `json:"temp_f"`
	IsDay        int       `json:"is_day"`
	Condition    Condition `json:"condition"`
	WindMph      float64   `json:"wind_mph"`
	WindKph      float64   `json:"wind_kph"`
	WindDegree   int       `json:"wind_degree"`
	WindDir      string    `json:"wind_dir"`
	PressureMb   float64   `json:"pressure_mb"`
	PressureIn   float64   `json:"pressure_in"`
	PrecipMm     float64   `json:"precip_mm"`
	PrecipIn     float64   `json:"precip_in"`
	Humidity     int       `json:"humidity"`
	Cloud        int       `json:"cloud"`
	FeelslikeC   float64   `json:"feelslike_c"`
	FeelslikeF   float64   `json:"feelslike_f"`
	WindchillC   float64   `json:"windchill_c"`
	WindchillF   float64   `json:"windchill_f"`
	HeatindexC   float64   `json:"heatindex_c"`
	HeatindexF   float64   `json:"heatindex_f"`
	DewpointC    float64   `json:"dewpoint_c"`
	DewpointF    float64   `json:"dewpoint_f"`
	WillItRain   int       `json:"will_it_rain"`
	ChanceOfRain string    `json:"chance_of_rain"`
	WillItSnow   int       `json:"will_it_snow"`
	ChanceOfSnow string    `json:"chance_of_snow"`
	VisKm        float64   `json:"vis_km"`
	VisMiles     float64   `json:"vis_miles"`
	GustMph      float64   `json:"gust_mph"`
	GustKph      float64   `json:"gust_kph"`
	UV           int       `json:"uv"`
}

type Condition struct {
	Text string `json:"text"`
	Icon string `json:"icon"`
	Code int    `json:"code"`
}
