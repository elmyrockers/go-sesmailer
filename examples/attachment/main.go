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
		SetSubject("Email With Attachments").
		SetBody("This is test email with attachments").
		Attach( "cat.webp", "" ). //<----------------------Add a few pictures as attachments
		Attach( "rabbit.jpg", "" ).
		Attach( "dog.jpg", "" ).
		Dump().
		Send()
}