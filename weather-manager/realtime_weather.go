package main

import "time"

type RealtimeWeather struct {
	Location struct {
		Name           string    `json:"name" db:"sName"`
		Region         string    `json:"region" db:"sRegion"`
		Country        string    `json:"country" db:"sCountry"`
		Latitude       float64   `json:"lat" db:"fLatitude"`
		Longitude      float64   `json:"lon" db:"fLongitude"`
		TimezoneID     string    `json:"tz_id" db:"sTimezoneID"`
		LocaltimeEpoch int64     `json:"localtime_epoch" db:"iLocaltimeEpoch"`
		Localtime      time.Time `json:"localtime" db:"dLocaltime"`
	} `json:"location"`
	Current struct {
		LastUpdatedEpoch int     `json:"last_updated_epoch" db:"iLastUpdatedEpoch"`
		LastUpdated      string  `json:"last_updated" db:"sLastUpdated"`
		TempC            float64 `json:"temp_c" db:"fTempC"`
		TempF            float64 `json:"temp_f" db:"fTempF"`
		IsDay            int     `json:"is_day" db:"iIsDay"`
		Condition        struct {
			Text string `json:"text" db:"sText"`
			Icon string `json:"icon" db:"sIcon"`
			Code int    `json:"code" db:"iCode"`
		} `json:"condition"`
		WindMph    float64 `json:"wind_mph" db:"fWindMph"`
		WindKph    float64 `json:"wind_kph" db:"fWindKph"`
		WindDegree int     `json:"wind_degree" db:"iWindDegree"`
		WindDir    string  `json:"wind_dir" db:"sWindDir"`
		PressureMb float64 `json:"pressure_mb" db:"fPressureMb"`
		PressureIn float64 `json:"pressure_in" db:"fPressureIn"`
		PrecipMm   float64 `json:"precip_mm" db:"fPrecipMm"`
		PrecipIn   float64 `json:"precip_in" db:"fPrecipIn"`
		Humidity   int     `json:"humidity" db:"iHumidity"`
		Cloud      int     `json:"cloud" db:"iCloud"`
		FeelslikeC float64 `json:"feelslike_c" db:"fFeelslikeC"`
		FeelslikeF float64 `json:"feelslike_f" db:"fFeelslikeF"`
		VisKm      float64 `json:"vis_km" db:"fVisKm"`
		VisMiles   float64 `json:"vis_miles" db:"fVisMiles"`
		UV         int     `json:"uv" db:"iUV"`
		GustMph    float64 `json:"gust_mph" db:"fGustMph"`
		GustKph    float64 `json:"gust_kph" db:"fGustKph"`
	} `json:"current"`
}
