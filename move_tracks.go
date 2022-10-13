package main

import (
	"context"
	"fmt"
	"log"

	"github.com/zmb3/spotify/v2"
)

func addTracksToContinuousPlaylist(client *spotify.Client) {
	const PLAYLIST_NAME = "Continuous Release Radar"

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
