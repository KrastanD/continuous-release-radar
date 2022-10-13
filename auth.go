package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

func authFunc() *spotify.Client {
	const redirectURI = "http://localhost:8080/callback"

	var auth = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopePlaylistModifyPrivate, spotifyauth.ScopePlaylistModifyPublic))
	var ch = make(chan *spotify.Client)
	var state = uuid.New().String()

	var fullCompleteAuth = handleAuthRedirect(auth, state, ch)

	http.HandleFunc("/callback", fullCompleteAuth)

	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	// wait for auth to complete
	client := <-ch
	return client
}

func handleAuthRedirect(auth *spotifyauth.Authenticator, state string, ch chan *spotify.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tok, err := auth.Token(r.Context(), state, r)
		if err != nil {
			http.Error(w, "Couldn't get token", http.StatusForbidden)
			log.Fatal(err)
		}
		if st := r.FormValue("state"); st != state {
			http.NotFound(w, r)
			log.Fatalf("State mismatch: %s != %s\n", st, state)
		}

		// use the token to get an authenticated client
		client := spotify.New(auth.Client(r.Context(), tok))
		fmt.Fprintf(w, "Login Completed!")
		ch <- client
	}
}
