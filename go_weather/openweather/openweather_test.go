package openweather

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchWeatherForCity(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("{\"id\": 12345, \"name\": \"Twin Peaks\"}"))
	}))
	defer server.Close()
	client := Client{baseUrl: server.URL, apiKey: "noKey", client: server.Client()}

	weatherForCity, _ := client.FetchWeatherForCity("Twin Peaks")

	if weatherForCity.CityName != "Twin Peaks" {
		t.Errorf("expected city to be 'Twin Peaks', but it was '%s'", weatherForCity.CityName)
	}
}
