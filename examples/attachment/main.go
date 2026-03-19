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
		AddAttachment( "cat.webp", "" ). //<----------------------Add a few pictures as attachments
		AddAttachment( "rabbit.jpg", "" ).
		AddAttachment( "dog.jpg", "" ).
		SetDebug(2).
		Send()
}