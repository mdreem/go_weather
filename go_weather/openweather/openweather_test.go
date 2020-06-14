package openweather

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchWeatherForCity(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		_, _ = rw.Write([]byte("{\"id\": 12345, \"name\": \"Twin Peaks\"}"))
	}))
	defer server.Close()
	client := Client{baseUrl: server.URL, apiKey: "noKey", client: server.Client()}

	weatherForCity, _ := client.FetchWeatherForCity("Twin Peaks")

	if weatherForCity.CityName != "Twin Peaks" {
		t.Errorf("expected city to be 'Twin Peaks', but it was '%s'", weatherForCity.CityName)
	}
}

func TestOpenWeatherFactoryMethod(t *testing.T) {
	client := CreateClient("someApiKey")

	if client.baseUrl != "https://api.openweathermap.org" {
		t.Errorf("baseUrl is '%s'", client.baseUrl)
	}

	if client.apiKey != "someApiKey" {
		t.Errorf("expected apiKey to be 'someApiKey' but is was '%s'", client.apiKey)
	}
}
