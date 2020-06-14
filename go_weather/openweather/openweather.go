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

func (o Client) FetchWeatherForCity(city string) (data.Weather, error) {
	request, err := o.createRequest(city)
	if err != nil {
		log.Printf("an error occured when creating the request: %v", err)
		return data.Weather{}, err
	}

	response, err := o.client.Do(request)
	if err != nil {
		log.Printf("an error occured when executing the request: %v", err)
		return data.Weather{}, err
	}

	weatherResponse, err := o.convertResponse(err, response)
	if err != nil {
		log.Printf("an error occured while converting the response: %v", err)
		return data.Weather{}, err
	}

	weather := data.Weather{
		Temperature: weatherResponse.Main.Temp,
		Pressure:    weatherResponse.Main.Pressure,
		CityName:    weatherResponse.Name,
		CityId:      weatherResponse.Id,
	}
	return weather, nil
}

func (o Client) createRequest(city string) (*http.Request, error) {
	request, err := http.NewRequest("GET", o.baseUrl+"/data/2.5/weather", nil)
	if err != nil {
		return nil, err
	}

	q := request.URL.Query()
	q.Add("q", city)
	q.Add("appid", o.apiKey)
	request.URL.RawQuery = q.Encode()
	return request, nil
}

func (o Client) convertResponse(err error, response *http.Response) (WeatherResponse, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return WeatherResponse{}, err
	}

	weatherResponse := WeatherResponse{}
	err = json.Unmarshal(body, &weatherResponse)
	if err != nil {
		return WeatherResponse{}, err
	}
	return weatherResponse, nil
}
