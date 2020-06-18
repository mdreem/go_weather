package endpoints

import (
	"../data"
	"errors"
	"github.com/stretchr/testify/assert"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type Dummy struct {
	err error
}

func (o Dummy) FetchWeatherForCity(_ string) (data.Weather, error) {
	weather := data.Weather{
		Temperature: 30.1,
		Pressure:    1000,
		CityName:    "Twin Peaks",
		CityId:      123,
	}
	return weather, o.err
}

func TestFetchWeatherForCity(t *testing.T) {
	assertion := assert.New(t)
	request := http.Request{}

	responseRecorder := httptest.NewRecorder()
	controller := WeatherDataController{OpenWeatherMapClient: Dummy{}}
	controller.CityHandler(responseRecorder, &request)

	assertion.Equal(200, responseRecorder.Code)
}

func TestFetchWeatherForCityRespondsWithError(t *testing.T) {
	assertion := assert.New(t)
	request := http.Request{}

	responseRecorder := httptest.NewRecorder()
	controller := WeatherDataController{OpenWeatherMapClient: Dummy{err: errors.New("blubb")}}
	controller.CityHandler(responseRecorder, &request)

	assertion.Equal(500, responseRecorder.Code)
}

func TestRespondNoData(t *testing.T) {
	assertion := assert.New(t)
	responseRecorder := httptest.NewRecorder()
	respond(responseRecorder, nil, 202)

	assertion.Equal(202, responseRecorder.Code)
}

func TestRespondWithData(t *testing.T) {
	assertion := assert.New(t)
	type TestData struct {
		SomeTest float64 `json:"someTest"`
	}

	responseRecorder := httptest.NewRecorder()
	respond(responseRecorder, TestData{SomeTest: 2.5}, 200)

	assertion.Equal(200, responseRecorder.Code)

	bodyString := strings.TrimSpace(responseRecorder.Body.String())
	assertion.Equal("{\"someTest\":2.5}", bodyString)
}

func TestRespondWithError(t *testing.T) {
	assertion := assert.New(t)
	responseRecorder := httptest.NewRecorder()
	respond(responseRecorder, math.Inf(1), 200)

	assertion.Equal(500, responseRecorder.Code)
}
