package openweather

import (
	"../data"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Client struct {
	baseUrl string
	apiKey  string
	client  *http.Client
}

func CreateClient(apiKey string) Client {
	return Client{baseUrl: "https://api.openweathermap.org", apiKey: apiKey, client: &http.Client{}}
}

func (o Client) FetchWeatherForCity(city string) data.Weather {
	request := o.createRequest(city)

	response, err := o.client.Do(request)
	if err != nil {
		log.Printf("an error occured: %v", err)
	}

	weatherResponse := o.convertResponse(err, response)

	weather := data.Weather{
		Temperature: weatherResponse.Main.Temp,
		Pressure:    weatherResponse.Main.Pressure,
		CityName:    weatherResponse.Name,
		CityId:      weatherResponse.Id,
	}
	return weather
}

func (o Client) createRequest(city string) *http.Request {
	request, err := http.NewRequest("GET", o.baseUrl+"/data/2.5/weather", nil)
	if err != nil {
		log.Printf("an error occured: %v", err)
	}

	q := request.URL.Query()
	q.Add("q", city)
	q.Add("appid", o.apiKey)
	request.URL.RawQuery = q.Encode()
	return request
}

func (o Client) convertResponse(err error, response *http.Response) WeatherResponse {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("an error occured: %v", err)
	}

	weatherResponse := WeatherResponse{}
	err = json.Unmarshal(body, &weatherResponse)
	if err != nil {
		log.Printf("an error occured: %v", err)
	}
	return weatherResponse
}
