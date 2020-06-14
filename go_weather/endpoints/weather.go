package endpoints

import (
	"../openweather"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type WeatherDataController struct {
	OpenWeatherMapClient openweather.Client
}

func (c WeatherDataController) Run() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/weather/{city}", c.CityHandler)
	r.Use(responseHeaderMiddleware)

	http.Handle("/", r)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Println("error occurred while starting http server:", err)
	}
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

	err = json.NewEncoder(writer).Encode(weather)
	if err != nil {
		log.Println("error occurred while converting data:", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func HomeHandler(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(200)
	_, err := fmt.Fprintf(writer, "Home\n")
	if err != nil {
		log.Println("error occurred:", err)
	}
}

func responseHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
