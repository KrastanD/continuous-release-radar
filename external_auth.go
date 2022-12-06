package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

func getExternalAuth() *spotify.Client {
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
	_, doesNTFYTopicExist := os.LookupEnv("NTFY_TOPIC")
	if !doesNTFYTopicExist {
		fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)
	} else {
		ntfy_topic := os.Getenv("NTFY_TOPIC")
		http.Post("https://ntfy.sh/"+ntfy_topic, "text/plain",
			strings.NewReader("Please log in to Spotify by visiting the following page in your browser:\n"+url))
	}

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

		saveTokenToFile(tok)

		// use the token to get an authenticated client
		client := spotify.New(auth.Client(r.Context(), tok))
		fmt.Fprintf(w, "Login Completed!")
		ch <- client
	}
}

func saveTokenToFile(tok *oauth2.Token) {
	marshaledToken, marshallingErr := json.Marshal(tok)
	if marshallingErr != nil {
		log.Print(marshallingErr)
	}
	writeErr := os.WriteFile("token.txt", marshaledToken, 0644)
	if writeErr != nil {
		log.Print(writeErr)
	}
}
