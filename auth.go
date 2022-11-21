package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

func getAuthClient() *spotify.Client {
	var tok *oauth2.Token
	//check for token in local storage
	data, readErr := os.ReadFile("token.txt")
	unMarshalErr := json.Unmarshal(data, &tok)
	if readErr != nil || unMarshalErr != nil {
		log.Print("Unable to read file or unmarshal json of token")
		//get a new token
		return getExternalAuth()
	}
	//use existing token to create an auth client
	var auth = spotifyauth.New(spotifyauth.WithRedirectURL(""), spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopePlaylistModifyPrivate, spotifyauth.ScopePlaylistModifyPublic))
	return spotify.New(auth.Client(context.TODO(), tok))
}
