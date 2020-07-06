package openweather

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchWeatherForCity(t *testing.T) {
	assertion := assert.New(t)

	server := createTestServer(200, t)
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

	server := createTestServer(200, t)
	defer server.Close()
	client := Client{baseUrl: string(0x7f), apiKey: "noKey", client: server.Client()}

	_, err := client.FetchWeatherForCity("Twin Peaks")

	assertion.Error(err)
	assertion.Regexp("invalid control character in URL", err.Error())
}

func TestFetchWeatherForCityInvalidTarget(t *testing.T) {
	assertion := assert.New(t)

	server := createTestServer(200, t)
	defer server.Close()
	client := Client{baseUrl: "http://nothing", apiKey: "noKey", client: server.Client()}

	_, err := client.FetchWeatherForCity("Twin Peaks")

	assertion.Error(err)
	assertion.Regexp("dial tcp: lookup nothing", err.Error())
}

func TestFetchWeatherForCityBadResponseCode(t *testing.T) {
	assertion := assert.New(t)

	t.Run("Bad response code", func(t *testing.T) {
		server := createTestServer(500, t)
		defer server.Close()
		client := Client{baseUrl: server.URL, apiKey: "noKey", client: server.Client()}

		_, err := client.FetchWeatherForCity("Twin Peaks")

		assertion.Error(err)
		assertion.Regexp("wrong response code received", err.Error())
	})
}

func TestFetchWeather(t *testing.T) {
	assertion := assert.New(t)

	t.Run("Reading error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Set("Content-Length", "1")
		}))
		defer server.Close()
		client := Client{baseUrl: server.URL, apiKey: "noKey", client: server.Client()}

		_, err := client.FetchWeatherForCity("Twin Peaks")

		assertion.EqualError(err, "unexpected EOF")
	})

	t.Run("Unmarshalling error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			_, err := rw.Write([]byte("<noJson></noJson>"))
			if err != nil {
				t.Fatal("write failed")
			}
		}))
		defer server.Close()
		client := Client{baseUrl: server.URL, apiKey: "noKey", client: server.Client()}

		_, err := client.FetchWeatherForCity("Twin Peaks")

		assertion.EqualError(err, "invalid character '<' looking for beginning of value")
	})
}

func createTestServer(statusCode int, t *testing.T) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(statusCode)
		_, err := rw.Write([]byte("{\"id\": 12345, \"name\": \"Twin Peaks\"}"))
		if err != nil {
			t.Fatal("write failed")
		}
	}))
	return server
}
