package main

import (
	"log"
	"os"
)

func checkForEnvVars() {
	_, doesIdExist := os.LookupEnv("SPOTIFY_ID")
	_, doesSecretExist := os.LookupEnv("SPOTIFY_SECRET")
	if !doesIdExist || !doesSecretExist {
		log.Fatal("Missing environment variables")
	}
}
