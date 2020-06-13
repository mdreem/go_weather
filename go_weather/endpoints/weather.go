package endpoints

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func Initialize() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/weather/{city}", CityHandler)
	r.Use(responseHeaderMiddleware)

	http.Handle("/", r)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println("Error occurred while starting http server:", err)
	}
}

func CityHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	writer.WriteHeader(http.StatusOK)

	weather := Weather{
		Temperature: 0,
		Pressure:    0,
		CityName:    vars["city"],
		CityId:      0,
	}
	err := json.NewEncoder(writer).Encode(weather)
	if err != nil {
		fmt.Println("Error occurred:", err)
	}
}

func HomeHandler(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(200)
	_, err := fmt.Fprintf(writer, "Home\n")
	if err != nil {
		fmt.Println("Error occurred:", err)
	}
}

func responseHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
