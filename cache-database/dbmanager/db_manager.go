package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoUri := "mongodb://cache-database-service:8092" 
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))

	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database("Weather").Collection("documents")
	createTTLIndex(collection)

	http.HandleFunc("/save-forecast", func(w http.ResponseWriter, r *http.Request) {
		saveCurrentWeather(w, r, collection)
	})

	http.HandleFunc("/get-forecast", func(w http.ResponseWriter, r *http.Request) {
		getCurrentWeather(w, r, collection)
	})

	http.ListenAndServe(":8090", nil)
}

func getLocation(loc string) ([]Location, error) {
	url := fmt.Sprintf("https://weatherapi-com.p.rapidapi.com/search.json?q=%s", loc)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", "2fdb8b6f67mshb1178afe669a190p1a700djsnc0d0e7bce537")
	req.Header.Add("X-RapidAPI-Host", "weatherapi-com.p.rapidapi.com")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	var locations []Location
	err = json.NewDecoder(res.Body).Decode(&locations)
	if err != nil {
		return nil, fmt.Errorf("decoding response failed: %w", err)
	}

	return locations, nil
}

func saveCurrentWeather(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	var item ForecastWeather

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Log the raw request body
	log.Printf("Raw request body: %s\n", string(bodyBytes))

	err = json.Unmarshal(bodyBytes, &item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	item.CreatedAt = time.Now()
	result, err := collection.InsertOne(ctx, item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Log the inserted ID
	log.Printf("Inserted ID: %v\n", result.InsertedID)

	fmt.Fprintf(w, "Item saved successfully in database. ")
}

func getCurrentWeather(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	// Get location parameter from query string
	locationName := r.URL.Query().Get("location")

	if locationName == "" {
		http.Error(w, "Missing 'location' parameter", http.StatusBadRequest)
		return
	}

	locations, err := getLocation(locationName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If no locations are found, return an error
	if len(locations) == 0 {
		http.Error(w, "No locations found", http.StatusNotFound)
		return
	}

	// Use the first location from the slice
	location := locations[0]

	// Define a filter for the search
	filter := bson.D{
		{"location.name", location.Name},
		{"location.region", location.Region},
	}

	fmt.Println("Location Name: ", location.Name)
	fmt.Println("Location Region: ", location.Region)

	var result ForecastWeather
	err = collection.FindOne(context.Background(), filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Return a 404 status code and error message if no document was found
			http.Error(w, "No document was found!", http.StatusNotFound)
		} else {
			// Return a 500 status code and error message for any other error
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			fmt.Println("FindOne error:", err)
		}
		return
	}

	// Marshal or convert user object back to json and write to response
	response, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func createTTLIndex(collection *mongo.Collection) {
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"createdAt": 1}, // Index key
		Options: options.Index().SetExpireAfterSeconds(300), // Expire after 5 minutes
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatal(err)
	}
}


type Location struct {
	Name           string  `bson:"name,omitempty" json:"name,omitempty"`
	Region         string  `bson:"region,omitempty" json:"region,omitempty"`
	Country        string  `bson:"country,omitempty" json:"country,omitempty"`
	Lat            float64 `bson:"lat,omitempty" json:"lat,omitempty"`
	Lon            float64 `bson:"lon,omitempty" json:"lon,omitempty"`
	TzId           string  `bson:"tz_id,omitempty" json:"tz_id,omitempty"`
	LocaltimeEpoch int64   `bson:"localtime_epoch,omitempty" json:"localtime_epoch,omitempty"`
	Localtime      string  `bson:"localtime,omitempty" json:"localtime,omitempty"`
}

type ForecastWeather struct {
	Location Location `bson:"location" json:"location"`
	Current  Current  `bson:"current" json:"current"`
	Forecast Forecast `bson:"forecast" json:"forecast"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

type Condition struct {
	Text string `bson:"text,omitempty" json:"text,omitempty"`
	Icon string `bson:"icon,omitempty" json:"icon,omitempty"`
	Code int    `bson:"code,omitempty" json:"code,omitempty"`
}

type Current struct {
	LastUpdatedEpoch int64     `bson:"last_updated_epoch,omitempty" json:"last_updated_epoch,omitempty"`
	LastUpdated      string    `bson:"last_updated,omitempty" json:"last_updated,omitempty"`
	TempC            float64   `bson:"temp_c,omitempty" json:"temp_c,omitempty"`
	TempF            float64   `bson:"temp_f,omitempty" json:"temp_f,omitempty"`
	IsDay            int       `bson:"is_day,omitempty" json:"is_day,omitempty"`
	Condition        Condition `bson:"condition,omitempty" json:"condition,omitempty"`
	WindMph          float64   `bson:"wind_mph,omitempty" json:"wind_mph,omitempty"`
	WindKph          float64   `bson:"wind_kph,omitempty" json:"wind_kph,omitempty"`
	WindDegree       int       `bson:"wind_degree,omitempty" json:"wind_degree,omitempty"`
	WindDir          string    `bson:"wind_dir,omitempty" json:"wind_dir,omitempty"`
	PressureMb       float64   `bson:"pressure_mb,omitempty" json:"pressure_mb,omitempty"`
	PressureIn       float64   `bson:"pressure_in,omitempty" json:"pressure_in,omitempty"`
	PrecipMm         float64   `bson:"precip_mm,omitempty" json:"precip_mm,omitempty"`
	PrecipIn         float64   `bson:"precip_in,omitempty" json:"precip_in,omitempty"`
	Humidity         int       `bson:"humidity,omitempty" json:"humidity,omitempty"`
	Cloud            int       `bson:"cloud,omitempty" json:"cloud,omitempty"`
	FeelslikeC       float64   `bson:"feelslike_c,omitempty" json:"feelslike_c,omitempty"`
	FeelslikeF       float64   `bson:"feelslike_f,omitempty" json:"feelslike_f,omitempty"`
	VisKm            float64   `bson:"vis_km,omitempty" json:"vis_km,omitempty"`
	VisMiles         float64   `bson:"vis_miles,omitempty" json:"vis_miles,omitempty"`
	UV               int       `bson:"uv,omitempty" json:"uv,omitempty"`
	GustMph          float64   `bson:"gust_mph,omitempty" json:"gust_mph,omitempty"`
	GustKph          float64   `bson:"gust_kph,omitempty" json:"gust_kph,omitempty"`
}

type Forecast struct {
	Forecastday []Forecastday `bson:"forecastday" json:"forecastday"`
}

type Forecastday struct {
	Date      string `bson:"date" json:"date"`
	DateEpoch int64  `bson:"date_epoch" json:"date_epoch"`
	Day       Day    `bson:"day" json:"day"`
}

type Day struct {
	MaxtempC          float64   `bson:"maxtemp_c" json:"maxtemp_c"`
	MaxtempF          float64   `bson:"maxtemp_f" json:"maxtemp_f"`
	MintempC          float64   `bson:"mintemp_c" json:"mintemp_c"`
	MintempF          float64   `bson:"mintemp_f" json:"mintemp_f"`
	AvgtempC          float64   `bson:"avgtemp_c" json:"avgtemp_c"`
	AvgtempF          float64   `bson:"avgtemp_f" json:"avgtemp_f"`
	MaxwindMph        float64   `bson:"maxwind_mph" json:"maxwind_mph"`
	MaxwindKph        float64   `bson:"maxwind_kph" json:"maxwind_kph"`
	TotalprecipMm     float64   `bson:"totalprecip_mm" json:"totalprecip_mm"`
	TotalprecipIn     float64   `bson:"totalprecip_in" json:"totalprecip_in"`
	TotalsnowCm       float64   `bson:"totalsnow_cm" json:"totalsnow_cm"`
	AvgvisKm          float64   `bson:"avgvis_km" json:"avgvis_km"`
	AvgvisMiles       float64   `bson:"avgvis_miles" json:"avgvis_miles"`
	Avghumidity       float64   `bson:"avghumidity" json:"avghumidity"`
	DailyWillItRain   int       `bson:"daily_will_it_rain" json:"daily_will_it_rain"`
	DailyChanceOfRain string    `bson:"daily_chance_of_rain" json:"daily_chance_of_rain"`
	DailyWillItSnow   int       `bson:"daily_will_it_snow" json:"daily_will_it_snow"`
	DailyChanceOfSnow string    `bson:"daily_chance_of_snow" json:"daily_chance_of_snow"`
	Condition         Condition `bson:"condition" json:"condition"`
	UV                float64   `bson:"uv" json:"uv"`
}
