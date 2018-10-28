package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	location := os.Args[1]
	if len(os.Args) == 1 || location == "" {
		log.Fatal("Please provide a location to get a weather reading")
		return
	}

	weatherData, err := getWeatherData(location)
	if err != nil {
		log.Fatal("😭 Failed to get the weather: " + err.Error())
	}

	report, err := formatWeather(weatherData)

	fmt.Printf("Right now in %s:\n%s\n", location, report)
}

func getWeatherData(location string) ([]byte, error) {
	appID := os.Getenv("OPENWEATHER_APP_ID")
	url := fmt.Sprintf(
		"http://api.openweathermap.org/data/2.5/weather?q=%s&APPID=%s",
		location,
		appID,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	return bytes, nil
}

func formatWeather(bytes []byte) (string, error) {
	var payload map[string]interface{}
	if err := json.Unmarshal(bytes, &payload); err != nil {
		return "", err
	}

	weatherSlice := payload["weather"].([]interface{})
	weather := weatherSlice[0].(map[string]interface{})
	weatherDescription := weather["description"].(string)

	main := payload["main"].(map[string]interface{})
	tempKelvin := main["temp"].(float64)
	tempCelsius := kelvinToCelsius(tempKelvin)

	weatherReport := fmt.Sprintf(
		"Temperature: %f°c\nThe weather situation: %s\n",
		tempCelsius,
		weatherDescription,
	)

	return weatherReport, nil
}

func kelvinToCelsius(kelvin float64) float64 {
	return kelvin - 273.15
}
