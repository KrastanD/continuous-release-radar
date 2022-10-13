package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/zmb3/spotify/v2"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	_, exists := os.LookupEnv("SPOTIFY_ID")
	var client *spotify.Client
	if exists {
		client = authFunc()
	}
	addTracksToContinuousPlaylist(client)
}
