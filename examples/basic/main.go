package main

import (
	"log"

	"github.com/elmyrockers/go-sesmailer"
)

func main() {
	// Create mailer
	sesmailer.New()


	log.Println("Test runs")
}