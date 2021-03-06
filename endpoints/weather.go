package endpoints

import (
	"../openweather"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type WeatherDataController struct {
	OpenWeatherMapClient openweather.WeatherFetcher
	KeycloakClientSecret string
}

func (c WeatherDataController) CityHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	city := vars["city"]
	log.Printf("fetching data for city '%s'", city)
	weather, err := c.OpenWeatherMapClient.FetchWeatherForCity(city)
	if err != nil {
		log.Println("error occurred while fetching weather data:", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	respond(writer, weather, http.StatusOK)
}

func (c WeatherDataController) HomeHandler(writer http.ResponseWriter, _ *http.Request) {
	type homeHandler struct {
		Home string
	}
	respond(writer, homeHandler{Home: "sweet Home"}, http.StatusOK)
}

func respond(writer http.ResponseWriter, data interface{}, statusCode int) {
	if data == nil {
		writer.WriteHeader(statusCode)
		return
	}

	err := json.NewEncoder(writer).Encode(data)
	if err != nil {
		log.Println("error occurred while converting data:", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(statusCode)
}

func responseHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
