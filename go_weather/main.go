package main

import (
	"./endpoints"
	"./openweather"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	configuration := loadConfiguration()

	openWeatherMapClient := openweather.Client{ApiKey: configuration.ApiToken}

	controller := endpoints.WeatherDataController{OpenWeatherMapClient: openWeatherMapClient}
	controller.Run()
}

type Configuration struct {
	ApiToken string
}

func loadConfiguration() Configuration {
	file, err := os.Open("configuration.json")
	if err != nil {
		fmt.Println("Error loading configuration file:", err)
		os.Exit(1)
	}

	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("Error closing file:", err)
			os.Exit(1)
		}
	}()

	decoder := json.NewDecoder(file)
	configuration := Configuration{}

	err = decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("Error decoding configuration file:", err)
	}
	return configuration
}
