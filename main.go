package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	_, exists := os.LookupEnv("SPOTIFY_ID")
	if exists {
		authFunc()
	}
}
