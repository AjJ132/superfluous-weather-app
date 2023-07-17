package main

import (
	"encoding/json"
	"io"
	"net/http"
)

func main() {

}

func getRealtimeWeather() *RealtimeWeather {
	url := "https://weatherapi-com.p.rapidapi.com/current.json?q=30643"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", "2fdb8b6f67mshb1178afe669a190p1a700djsnc0d0e7bce537")
	req.Header.Add("X-RapidAPI-Host", "weatherapi-com.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var weather RealtimeWeather
	_ = json.Unmarshal(body, &weather)

	return &weather
}
