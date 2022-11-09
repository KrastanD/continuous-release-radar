package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

func getTokenFromLocalStorage(data []byte) *spotify.Client {
	var client *spotify.Client
	var tok *oauth2.Token
	err := json.Unmarshal(data, &tok)
	if err != nil {
		log.Print("Unable to unmarshal json of token")
		client = getExternalAuth()
	} else {
		const redirectURI = "http://localhost:8080/callback"
		var auth = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopePlaylistModifyPrivate, spotifyauth.ScopePlaylistModifyPublic))
		client = spotify.New(auth.Client(context.TODO(), tok))
	}

	return client
}
