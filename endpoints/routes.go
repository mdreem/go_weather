package endpoints

import (
	"context"
	"github.com/coreos/go-oidc"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"strings"
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

type OAuth2Config struct {
	config *oauth2.Config
}

func initializeOpenIdConnect(clientSecret string) *auth {
	clientID := "weather"
	redirectUrl := "http://localhost:8000/oidc/callback"

	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, "http://localhost:8080/auth/realms/Weather")
	if err != nil {
		panic(err)
	}

	oauth2Config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectUrl,

		Endpoint: provider.Endpoint(),
		Scopes:   []string{oidc.ScopeOpenID, "profile", "email"},
	}

	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}

	var verifier = provider.Verifier(oidcConfig)

	auth := auth{
		oauth2Config: OAuth2Config{config: &oauth2Config},
		verifier:     verifier,
		ctx:          &ctx,
	}
	return &auth
}

func (auth auth) authenticationHandlerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth.authenticationHandler(next, w, r)
	})
}

func (auth auth) authenticationHandler(next http.Handler, writer http.ResponseWriter, request *http.Request) {
	rawAccessToken := request.Header.Get("Authorization")

	if rawAccessToken == "" {
		log.Printf("Token missing.")
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	parts := strings.Split(rawAccessToken, " ")
	if len(parts) != 2 {
		log.Printf("Token malformed.")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err := auth.verifier.Verify(*auth.ctx, parts[1])
	if err != nil {
		log.Printf("Token could not be verified. %v", err)
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	next.ServeHTTP(writer, request)
}
