package openweather

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchWeatherForCity(t *testing.T) {
	assertion := assert.New(t)

	server := createTestServer(200)
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

func TestFetchWeatherForCityInvalidUrl(t *testing.T) {
	assertion := assert.New(t)

	server := createTestServer(200)
	defer server.Close()
	client := Client{baseUrl: string(0x7f), apiKey: "noKey", client: server.Client()}

	_, err := client.FetchWeatherForCity("Twin Peaks")

	assertion.EqualError(err, "parse \u007f/data/2.5/weather: net/url: invalid control character in URL")
}

func TestFetchWeatherForCityInvalidTarget(t *testing.T) {
	assertion := assert.New(t)

	server := createTestServer(200)
	defer server.Close()
	client := Client{baseUrl: "http://nothing", apiKey: "noKey", client: server.Client()}

	_, err := client.FetchWeatherForCity("Twin Peaks")

	assertion.EqualError(err, "Get http://nothing/data/2.5/weather?appid=noKey&q=Twin+Peaks: dial tcp: lookup nothing: no such host")
}

func TestFetchWeatherForCityBadResponseCode(t *testing.T) {
	assertion := assert.New(t)

	server := createTestServer(500)
	defer server.Close()
	client := Client{baseUrl: server.URL, apiKey: "noKey", client: server.Client()}

	_, err := client.FetchWeatherForCity("Twin Peaks")

	assertion.EqualError(err, "wrong response code received")
}

func TestFetchWeatherForCityReadingError(t *testing.T) {
	assertion := assert.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Length", "1")
	}))
	defer server.Close()
	client := Client{baseUrl: server.URL, apiKey: "noKey", client: server.Client()}

	_, err := client.FetchWeatherForCity("Twin Peaks")

	assertion.EqualError(err, "unexpected EOF")
}

func TestFetchWeatherForCityUnmarshallingError(t *testing.T) {
	assertion := assert.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		_, _ = rw.Write([]byte("<noJson></noJson>"))
	}))
	defer server.Close()
	client := Client{baseUrl: server.URL, apiKey: "noKey", client: server.Client()}

	_, err := client.FetchWeatherForCity("Twin Peaks")

	assertion.EqualError(err, "invalid character '<' looking for beginning of value")
}

func createTestServer(statusCode int) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(statusCode)
		_, _ = rw.Write([]byte("{\"id\": 12345, \"name\": \"Twin Peaks\"}"))
	}))
	return server
}
