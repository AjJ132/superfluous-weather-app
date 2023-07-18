package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	fmt.Println("Starting the application...")
	handler := &Handler{
		Location: "30643",
		//databse will go here
	}

	http.HandleFunc("/realtime-weather", func(w http.ResponseWriter, r *http.Request) {
		handler.getRealtimeWeather(w, r)
	});

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil));


}

func (h * Handler) getRealtimeWeather(w http.ResponseWriter, r *http.Request) {
	//Decicde between cahce or not

	//first check database cache

	//if cache is avaliable, return cache

	//else call API

	//Use cache service to store cache and return data to user symultaneously
}


type Handler struct {
	//location variable from HTTP request
	Location string `json:"location"`
}