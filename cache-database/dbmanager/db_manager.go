package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoUri := os.Getenv("MONGODB_URI")
	if mongoUri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable.")
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))

	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database("Weather").Collection("documents")

	http.HandleFunc("/save-current-weather", func(w http.ResponseWriter, r *http.Request) {
		saveCurrentWeather(w, r, collection)
	})

	http.ListenAndServe(":8080", nil)
}

func saveCurrentWeather(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	var item Item

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Log the raw request body
	log.Printf("Raw request body: %s\n", string(bodyBytes))

	err = json.Unmarshal(bodyBytes, &item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	result, err := collection.InsertOne(ctx, item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Log the inserted ID
	log.Printf("Inserted ID: %v\n", result.InsertedID)

	fmt.Fprintf(w, "Item saved successfully in database. ")
}

type Item struct {
	Location Location `bson:"location"`
	Current  Current  `bson:"current"`
}

type Location struct {
	Name           string  `bson:"name"`
	Region         string  `bson:"region"`
	Country        string  `bson:"country"`
	Lat            float64 `bson:"lat"`
	Lon            float64 `bson:"lon"`
	TzID           string  `bson:"tz_id"`
	LocaltimeEpoch int64   `bson:"localtime_epoch"`
	Localtime      string  `bson:"localtime"`
}

type Condition struct {
	Text string `bson:"text"`
	Icon string `bson:"icon"`
	Code int    `bson:"code"`
}

type Current struct {
	LastUpdatedEpoch int64     `bson:"last_updated_epoch"`
	LastUpdated      string    `bson:"last_updated"`
	TempC            float64   `bson:"temp_c"`
	TempF            float64   `bson:"temp_f"`
	IsDay            int       `bson:"is_day"`
	Condition        Condition `bson:"condition"`
	WindMph          float64   `bson:"wind_mph"`
	WindKph          float64   `bson:"wind_kph"`
	WindDegree       int       `bson:"wind_degree"`
	WindDir          string    `bson:"wind_dir"`
	PressureMb       float64   `bson:"pressure_mb"`
	PressureIn       float64   `bson:"pressure_in"`
	PrecipMm         float64   `bson:"precip_mm"`
	PrecipIn         float64   `bson:"precip_in"`
	Humidity         int       `bson:"humidity"`
	Cloud            int       `bson:"cloud"`
	FeelslikeC       float64   `bson:"feelslike_c"`
	FeelslikeF       float64   `bson:"feelslike_f"`
	VisKm            float64   `bson:"vis_km"`
	VisMiles         float64   `bson:"vis_miles"`
	UV               int       `bson:"uv"`
	GustMph          float64   `bson:"gust_mph"`
	GustKph          float64   `bson:"gust_kph"`
}
