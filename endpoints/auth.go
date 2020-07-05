package endpoints

import (
	"context"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"strings"
)

type auth struct {
	oauth2Config authConfig
	verifier     *oidc.IDTokenVerifier
	ctx          *context.Context
}

type authConfig interface {
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
	GetClientID() string
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

func (config OAuth2Config) Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	return config.config.Exchange(ctx, code, opts...)
}

func (config OAuth2Config) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	return config.config.AuthCodeURL(state, opts...)
}

func (config OAuth2Config) GetClientID() string {
	return config.config.ClientID
}

func (auth auth) handleOAuth2Callback(writer http.ResponseWriter, request *http.Request) {
	oauth2Token, err := auth.oauth2Config.Exchange(*auth.ctx, request.URL.Query().Get("code"))
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not exchange code for token: %v", err)
		return
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Printf("missing token")
		return
	}

	_, err = writer.Write([]byte(rawIDToken))
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Printf("not able to write token")
		return
	}
}

func (auth auth) handleLogin(writer http.ResponseWriter, request *http.Request) {
	log.Printf("Data: %s", auth.oauth2Config.GetClientID())
	http.Redirect(writer, request, auth.oauth2Config.AuthCodeURL("main"), http.StatusFound)
}
