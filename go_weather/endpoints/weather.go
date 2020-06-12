package endpoints

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func Initialize() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/weather/{city}", CityHandler)

	http.Handle("/", r)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println("Error occurred while starting http server:", err)
	}
}

func CityHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	writer.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintf(writer, "City: %v\n", vars["city"])
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
