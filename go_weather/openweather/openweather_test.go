package openweather

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchWeatherForCity(t *testing.T) {
	assertion := assert.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		_, _ = rw.Write([]byte("{\"id\": 12345, \"name\": \"Twin Peaks\"}"))
	}))
	defer server.Close()
	client := Client{baseUrl: server.URL, apiKey: "noKey", client: server.Client()}

	weatherForCity, _ := client.FetchWeatherForCity("Twin Peaks")

	assertion.Equal("Twin Peaks", weatherForCity.CityName)
}

func TestOpenWeatherFactoryMethod(t *testing.T) {
	assertion := assert.New(t)
	client := CreateClient("someApiKey")

	assertion.Equal("https://api.openweathermap.org", client.baseUrl)
	assertion.Equal("someApiKey", client.apiKey)
}
