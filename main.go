package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type WeatherResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  string  `json:"timezone"`
	Currently struct {
		Icon                string  `json:"icon"`
		Time                int     `json:"time"`
		Summary             string  `json:"summary"`
		Temperature         float64 `json:"temperature"`
		ApparentTemperature float64 `json:"apparentTemperature"`
		PrecipProbability   float64 `json:"precipProbability"`
		WindSpeed           float64 `json:"windSpeed"`
		Humidity            float64 `json:"humidity"`
	} `json:"currently"`
	Daily struct {
		Data []struct {
			MoonPhase   float64 `json:"moonPhase"`
			SunriseTime int     `json:"sunriseTime"`
			SunsetTime  int     `json:"sunsetTime"`
		} `json:"data"`
	} `json:"daily"`
}

func getFloatFromEnv(envVar string, defaultValue float64) float64 {
	val := os.Getenv(envVar)
	if val == "" {
		return defaultValue
	}
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return defaultValue
	}
	return f
}

func getMoonPhaseEmoji(moonPhase float64) string {
	switch {
	case moonPhase < 0.125:
		return "🌑" // new moon
	case moonPhase < 0.25:
		return "🌒" // waxing crescent
	case moonPhase < 0.375:
		return "🌓" // first quarter
	case moonPhase < 0.5:
		return "🌔" // waxing gibbous
	case moonPhase < 0.625:
		return "🌕" // full moon
	case moonPhase < 0.75:
		return "🌖" // waning gibbous
	case moonPhase < 0.875:
		return "🌗" // last quarter
	default:
		return "🌘" // waning crescent
	}
}

func formatTime(timestamp int, timezone string) string {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return "N/A"
	}
	t := time.Unix(int64(timestamp), 0).In(loc)
	return t.Format("3:04 PM")
}

func main() {
	defaultLat := getFloatFromEnv("PIRATE_WEATHER_LAT", 40.7128)  // NYC latitude
	defaultLon := getFloatFromEnv("PIRATE_WEATHER_LON", -74.0060) // NYC longitude
	defaultUnits := os.Getenv("PIRATE_WEATHER_UNITS")
	if defaultUnits == "" {
		defaultUnits = "us"
	}

	lat := flag.Float64("lat", defaultLat, "Latitude (can also use PIRATE_WEATHER_LAT environment variable)")
	lon := flag.Float64("lon", defaultLon, "Longitude (can also use PIRATE_WEATHER_LON environment variable)")
	units := flag.String("units", defaultUnits, "Units system (us, si, ca, uk) (can also use PIRATE_WEATHER_UNITS environment variable)")
	flag.Parse()

	apiKey := os.Getenv("PIRATE_WEATHER_API_KEY")
	if apiKey == "" {
		fmt.Println("Please set PIRATE_WEATHER_API_KEY environment variable")
		return
	}

	validUnits := map[string]bool{"us": true, "si": true, "ca": true, "uk": true}
	if !validUnits[*units] {
		fmt.Println("Invalid units. Must be one of: us, si, ca, uk")
		return
	}

	url := fmt.Sprintf("https://api.pirateweather.net/forecast/%s/%.6f,%.6f?units=%s&exclude=minutely,hourly,alerts",
		apiKey, *lat, *lon, *units)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("API returned error: %s\n%s\n", resp.Status, string(body))
		return
	}

	var weather WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
		return
	}

	emoji := map[string]string{
		"clear-day":           "☀️",
		"clear-night":         getMoonPhaseEmoji(weather.Daily.Data[0].MoonPhase),
		"rain":                "🌧️",
		"snow":                "🌨️",
		"sleet":               "🌨️",
		"wind":                "🌬️",
		"fog":                 "🌫️",
		"cloudy":              "☁️",
		"partly-cloudy-day":   "🌤️",
		"partly-cloudy-night": "☁️",
		"thunderstorm":        "⛈️",
		"hail":                "🌨️",
		"none":                "🏴‍☠️",
	}

	icon := emoji[weather.Currently.Icon]
	if icon == "" {
		icon = "🏴‍☠️"
	}

	tempUnit := "°F"
	if *units == "si" || *units == "ca" || *units == "uk" {
		tempUnit = "°C"
	}

	windUnit := "mph"
	if *units == "si" {
		windUnit = "m/s"
	} else if *units == "ca" {
		windUnit = "km/h"
	}

	sunrise := formatTime(weather.Daily.Data[0].SunriseTime, weather.Timezone)
	sunset := formatTime(weather.Daily.Data[0].SunsetTime, weather.Timezone)

	fmt.Printf("Pirate Weather\n")
	fmt.Printf("📍 %.6f,%.6f\n", *lat, *lon)
	// fmt.Printf("🕰️ %s\n", weather.Timezone)
	fmt.Printf("%s %s\n", icon, weather.Currently.Summary)
	fmt.Printf("🌅 Sunrise:        %s\n", sunrise)
	fmt.Printf("🌇 Sunset:         %s\n", sunset)
	fmt.Printf("🌡️ Temperature:    %.1f%s\n", weather.Currently.Temperature, tempUnit)
	fmt.Printf("🌡️ Feels Like:     %.1f%s\n", weather.Currently.ApparentTemperature, tempUnit)
	fmt.Printf("☔️ Precip Chance:  %.0f%%\n", weather.Currently.PrecipProbability*100)
	fmt.Printf("💧 Humidity:       %.0f%%\n", weather.Currently.Humidity*100)
	fmt.Printf("💨 Wind Speed:     %.1f %s\n", weather.Currently.WindSpeed, windUnit)
}
