package openweather

import (
	"../endpoints"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Client struct {
	ApiKey string
}

func (o Client) FetchWeatherForCity(city string) endpoints.Weather {
	client := &http.Client{}
	request, err := http.NewRequest("GET", "https://api.openweathermap.org/data/2.5/weather", nil)
	if err != nil {
		log.Printf("an error occured: %v", err)
	}

	q := request.URL.Query()
	q.Add("q", city)
	q.Add("appid", o.ApiKey)
	request.URL.RawQuery = q.Encode()

	response, err := client.Do(request)
	if err != nil {
		log.Printf("an error occured: %v", err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("an error occured: %v", err)
	}

	weatherResponse := WeatherResponse{}
	err = json.Unmarshal(body, &weatherResponse)
	if err != nil {
		log.Printf("an error occured: %v", err)
	}

	weather := endpoints.Weather{
		Temperature: weatherResponse.Main.Temp,
		Pressure:    weatherResponse.Main.Pressure,
		CityName:    weatherResponse.Name,
		CityId:      weatherResponse.Id,
	}
	return weather
}
