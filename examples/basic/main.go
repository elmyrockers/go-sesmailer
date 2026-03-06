package main

import (
	// "github.com/davecgh/go-spew/spew"

	"context"
	_ "github.com/joho/godotenv/autoload"
	"github.com/elmyrockers/go-sesmailer"
)

func main() {
	// Create mailer
		sesmailer.New().
			SetFrom( "no-reply@xeno.com.my", "Xeno System" ).
			AddAddress( "elmyrockers@gmail.com", "Helmi Aziz" ).
			SetSubject( "Test Subjek 7" ).
			SetBody( "test sahaja body" ).
			SetAltBody( "test sahaja alt body" ).
			// SetDebug( 2 ).
			Send( context.Background() )
}


