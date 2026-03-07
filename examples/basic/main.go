package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/elmyrockers/go-sesmailer"
)

func main() {
	// Create mailer
		sesmailer.New().
			SetFrom( "no-reply@xeno.com.my", "Xeno System" ).
			AddAddress( "elmyrockers@gmail.com", "Helmi Aziz" ).
			AddReplyTo( "elmyrockers2@gmail.com", "Helmi Aziz 2" ).
			SetSubject( "Test subject" ).
			SetBody( "test body" ).
			SetAltBody( "test alt body" ).
			SetDebug( 2 ).
			Send()
}


