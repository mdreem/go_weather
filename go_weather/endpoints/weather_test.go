package endpoints

import (
	"../data"
	"errors"
	"log"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Dummy struct {
	err error
}

func (o Dummy) FetchWeatherForCity(city string) (data.Weather, error) {
	weather := data.Weather{
		Temperature: 30.1,
		Pressure:    1000,
		CityName:    "Twin Peaks",
		CityId:      123,
	}
	return weather, o.err
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

func TestFetchWeatherForCityRespondsWithError(t *testing.T) {
	request := http.Request{}

	responseRecorder := httptest.NewRecorder()
	controller := WeatherDataController{OpenWeatherMapClient: Dummy{err: errors.New("blubb")}}
	controller.CityHandler(responseRecorder, &request)

	if responseRecorder.Code != 500 {
		t.Errorf("expected code to be '500', but it was '%d'", responseRecorder.Code)
	}

	log.Printf("Response: '%s'", responseRecorder.Body.String())
}

func TestRespondNoData(t *testing.T) {
	responseRecorder := httptest.NewRecorder()
	respond(responseRecorder, nil, 200)

	if responseRecorder.Code != 200 {
		t.Errorf("expected code to be '200', but it was '%d'", responseRecorder.Code)
	}
}

func TestRespondWithData(t *testing.T) {
	type TestData struct {
		SomeTest float64 `json:"someTest"`
	}

	responseRecorder := httptest.NewRecorder()
	respond(responseRecorder, TestData{SomeTest: 2.5}, 200)

	if responseRecorder.Code != 200 {
		t.Errorf("expected code to be '200', but it was '%d'", responseRecorder.Code)
	}
}

func TestRespondWithError(t *testing.T) {
	responseRecorder := httptest.NewRecorder()
	respond(responseRecorder, math.Inf(1), 200)

	if responseRecorder.Code != 500 {
		t.Errorf("expected code to be '500', but it was '%d'", responseRecorder.Code)
	}
}
