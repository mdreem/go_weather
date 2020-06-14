package main

import (
	"./endpoints"
	"./openweather"
	"encoding/json"
	"log"
	"os"
)

func main() {
	configuration := loadConfiguration()

	openWeatherMapClient := openweather.CreateClient(configuration.ApiToken)

	controller := endpoints.WeatherDataController{OpenWeatherMapClient: openWeatherMapClient}
	controller.Run()
}

type Configuration struct {
	ApiToken string
}

func loadConfiguration() Configuration {
	file, err := os.Open("configuration.json")
	if err != nil {
		log.Println("error loading configuration file:", err)
		os.Exit(1)
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Println("error closing file:", err)
			os.Exit(1)
		}
	}()

	decoder := json.NewDecoder(file)
	configuration := Configuration{}

	err = decoder.Decode(&configuration)
	if err != nil {
		log.Println("error decoding configuration file:", err)
	}
	return configuration
}
