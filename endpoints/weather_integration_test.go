package endpoints

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type TokenResponse struct {
	IdToken string `json:"id_token"`
}

func getToken(t *testing.T) string {
	data := url.Values{}
	data.Set("client_id", "weather")
	data.Set("grant_type", "password")
	data.Set("client_secret", "01c2497e-1cfb-45bb-844d-482f75bd9a6a")
	data.Set("scope", "openid")
	data.Set("username", "testuser")
	data.Set("password", "password")

	const keycloakTokenUrl = "http://localhost:8080/auth/realms/Weather/protocol/openid-connect/token"
	request, err := http.NewRequest("POST", keycloakTokenUrl, strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{}
	response, _ := client.Do(request)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	tokenResponse := TokenResponse{}
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	return tokenResponse.IdToken
}

func TestFetchWeatherForCity(t *testing.T) {
	assertion := assert.New(t)

	request, err := http.NewRequest("GET", "/weather/NoCity", nil)
	if err != nil {
		t.Fatal(err)
	}
	bearerToken := fmt.Sprintf("Bearer %s", getToken(t))
	request.Header.Add("Authorization", bearerToken)

	controller := WeatherDataController{OpenWeatherMapClient: Dummy{}}
	handler := controller.SetupRoutes()

	responseRecorder := httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, request)

	assertion.Equal(200, responseRecorder.Code)
	assertion.Equal(ApplicationJson, responseRecorder.Header().Get("Content-Type"))
	assertion.Equal("{\"temperature\":30.1,\"pressure\":1000,\"cityName\":\"Twin Peaks\",\"cityId\":123}\n", responseRecorder.Body.String())
}

func TestHomeHandler(t *testing.T) {
	assertion := assert.New(t)

	request, err := http.NewRequest("GET", "/weather/", nil)
	if err != nil {
		t.Fatal(err)
	}
	bearerToken := fmt.Sprintf("Bearer %s", getToken(t))
	request.Header.Add("Authorization", bearerToken)

	controller := WeatherDataController{OpenWeatherMapClient: Dummy{}}
	handler := controller.SetupRoutes()

	responseRecorder := httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, request)

	assertion.Equal(200, responseRecorder.Code)
	assertion.Equal(ApplicationJson, responseRecorder.Header().Get("Content-Type"))
	assertion.Equal("{\"Home\":\"sweet Home\"}\n", responseRecorder.Body.String())
}
