package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

func authFunc() {
	const PLAYLIST_NAME = "Continuous Release Radar"
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

	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	userPlaylists, err := client.CurrentUsersPlaylists(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var crr_playlist *spotify.SimplePlaylist
	for _, p := range userPlaylists.Playlists {
		playlist := p
		if playlist.Name == PLAYLIST_NAME {
			crr_playlist = &playlist
			break
		}
	}

	var crr_search_result *spotify.SearchResult
	if crr_playlist == nil {
		crr_search_result, err = client.Search(context.Background(), PLAYLIST_NAME, spotify.SearchTypePlaylist)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Printf("CRR %v\n", crr_playlist.Name)
	}

	if crr_search_result.Playlists.Playlists[0].Name == PLAYLIST_NAME {
		crr_playlist = &crr_search_result.Playlists.Playlists[0]
	}

	if crr_playlist == nil {
		new_playlist, err := client.CreatePlaylistForUser(context.Background(), user.ID, PLAYLIST_NAME, "", false, false)
		if err != nil {
			log.Fatal(err)
		}
		crr_playlist = &(new_playlist.SimplePlaylist)
		fmt.Println("CRR Playlist created")
	} else {
		fmt.Println("CRR Playlist found")
	}

	search_result, err := client.Search(context.Background(), "Release Radar", spotify.SearchTypePlaylist)

	if err != nil {
		log.Fatal(err)
	}

	if search_result != nil && search_result.Playlists != nil {
		for _, playlist := range search_result.Playlists.Playlists {
			if playlist.Name == "Release Radar" {
				tracks, _ := client.GetPlaylistItems(context.Background(), playlist.ID)
				for _, songs := range tracks.Items {
					client.AddTracksToPlaylist(context.Background(), crr_playlist.ID, songs.Track.Track.ID)
				}
				break
			}
		}
	}

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
