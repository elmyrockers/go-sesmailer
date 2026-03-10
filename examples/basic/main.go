package main

import (
	"github.com/elmyrockers/go-sesmailer"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// Create mailer
	sesmailer.New().
		SetFrom("no-reply@xeno.com.my", "Xeno System").
		AddAddress("elmyrockers@gmail.com", "Helmi Aziz").
		AddReplyTo("elmyrockers2@gmail.com", "Helmi Aziz 2").
		SetSubject("Test subject 20").
		SetBody("test body").
		SetAltBody("test alt body").
		SetDebug(2).
		Send()
}
