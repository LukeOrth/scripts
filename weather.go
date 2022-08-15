package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
)

type WeatherResult struct {
	Weather     []*Weather
	Temperature *Temperature `json:"main"`
}

type Weather struct {
	Main string
}

type Temperature struct {
	Temp float64
}

func main() {
	apiKey := os.Getenv("OWM_API_KEY")
	city := os.Getenv("CURRENT_CITY")

	weather, err := getWeather(apiKey, city)
	if err != nil {
		log.Fatal(err)
	}

	temp, condition := math.Round(weather.Temperature.Temp), string(weather.Weather[0].Main)

	fmt.Printf("%0.f° %s\n", temp, getIcon(condition))
}

func getWeather(key string, city string) (*WeatherResult, error) {
	params := url.Values{
		"q":     []string{city},
		"appid": []string{key},
		"units": []string{"imperial"},
	}

	u := &url.URL{
		Scheme:   "https",
		Host:     "api.openweathermap.org",
		Path:     "/data/2.5/weather",
		RawQuery: params.Encode(),
	}

	res, err := http.Get(u.String())

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OWM weather query failed: %s", err)
	}

	defer res.Body.Close()

	var result WeatherResult

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("JSON decoding failed: %s", err)
	}

	return &result, nil
}

func getIcon(condition string) string {
	icon := map[string]string{
		"Thunderstorm": "",
		"Drizzle":      "",
		"Rain":         "",
		"Snow":         "",
		"Atmosphere":   "",
		"Clear":        "盛",
		"Clouds":       "",
	}

	return icon[condition]
}
