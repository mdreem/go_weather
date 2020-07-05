package endpoints

import (
	"github.com/gorilla/mux"
)

func (c WeatherDataController) SetupRoutes() *mux.Router {
	auth := initializeOpenIdConnect(c.KeycloakClientSecret)

	r := mux.NewRouter()

	weather := r.PathPrefix("/weather").Subrouter()
	weather.HandleFunc("/", c.HomeHandler).Methods("GET")
	weather.HandleFunc("/{city}", c.CityHandler).Methods("GET")
	weather.Use(responseHeaderMiddleware)
	weather.Use(auth.authenticationHandlerMiddleware)

	oidcHandler := r.PathPrefix("/oidc").Subrouter()
	oidcHandler.HandleFunc("/callback", auth.handleOAuth2Callback).Methods("GET")
	oidcHandler.HandleFunc("/login", auth.handleLogin).Methods("GET")

	return r
}
