package endpoints

import (
	"../data"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Dummy struct {
}

func (o Dummy) FetchWeatherForCity(city string) (data.Weather, error) {
	weather := data.Weather{
		Temperature: 30.1,
		Pressure:    1000,
		CityName:    "Twin Peaks",
		CityId:      123,
	}
	return weather, nil
}

func TestFetchWeatherForCity(t *testing.T) {
	request := http.Request{}

	responseRecorder := httptest.NewRecorder()
	controller := WeatherDataController{OpenWeatherMapClient: Dummy{}}
	controller.CityHandler(responseRecorder, &request)

	if responseRecorder.Code != 200 {
		t.Errorf("expected code to be '200', but it was '%d'", responseRecorder.Code)
	}

	log.Printf("Response: '%s'", responseRecorder.Body.String())
}
