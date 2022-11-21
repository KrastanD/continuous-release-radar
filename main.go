package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/zmb3/spotify/v2"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
	checkForEnvVars()
}

func main() {
	var client *spotify.Client = getAuthClient()
	addTracksToContinuousPlaylist(client)
	log.Print("Success!")
}
