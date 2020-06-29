package endpoints

import (
	"log"
	"net/http"
)

func (auth auth) handleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v", auth.ctx)
	oauth2Token, err := auth.oauth2Config.Exchange(*auth.ctx, r.URL.Query().Get("code"))
	if err != nil {
		panic(err)
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		panic("Missing token")
	}

	_, err = w.Write([]byte(rawIDToken))
	if err != nil {
		panic(err)
	}
}

func (auth auth) handleLogin(w http.ResponseWriter, r *http.Request) {
	log.Printf("Data: %s", auth.oauth2Config.ClientID)
	http.Redirect(w, r, auth.oauth2Config.AuthCodeURL("main"), http.StatusFound)
}
