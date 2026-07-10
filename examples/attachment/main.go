package main

import (
	"os"
	"log"
	"github.com/elmyrockers/go-sesmailer"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// Get image data in bytes
		catImg, err := os.ReadFile("cat.webp")
			if err != nil { log.Fatal(err) }
		rabbitImg, err := os.ReadFile("rabbit.jpg")
			if err != nil { log.Fatal(err) }
		dogImg, err := os.ReadFile("dog.jpg")
			if err != nil { log.Fatal(err) }

	// Create mailer
	sesmailer.New().
		SetFrom("no-reply@xeno.com.my", "Xeno System").
		AddAddress("elmyrockers@gmail.com", "Helmi Aziz").
		SetSubject("Email With Attachments").
		SetBody("This is test email with attachments").
		Attach( "cat.webp", catImg ). //<----------------------Add a few pictures as attachments
		Attach( "rabbit.jpg", rabbitImg ).
		Attach( "dog.jpg", dogImg ).
		Dump().
		Send()
}