package endpoints

import (
	"context"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
	"log"
	"net/http"
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

func (config OAuth2Config) Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	return config.config.Exchange(ctx, code, opts...)
}

func (config OAuth2Config) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	return config.config.AuthCodeURL(state, opts...)
}

func (config OAuth2Config) GetClientID() string {
	return config.config.ClientID
}

func (auth auth) handleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
	oauth2Token, err := auth.oauth2Config.Exchange(*auth.ctx, r.URL.Query().Get("code"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		panic("Missing token")
	}

	_, err = w.Write([]byte(rawIDToken))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
}

func (auth auth) handleLogin(writer http.ResponseWriter, request *http.Request) {
	log.Printf("Data: %s", auth.oauth2Config.GetClientID())
	http.Redirect(writer, request, auth.oauth2Config.AuthCodeURL("main"), http.StatusFound)
}
