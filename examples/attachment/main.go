package main

import (
	"os"
	"log"
	"github.com/elmyrockers/go-sesmailer"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// // Get image data in bytes
		catImg, err := os.ReadFile("cat.webp")
			if err != nil {
				log.Fatalf("Failed to read attachment: %v", err)
				return
			}

	// Create mailer
		mailer := sesmailer.New()

		messageID, err := mailer.
								SetFrom("noreply@xeno.com.my", "Xeno System").
								AddAddress("elmyrockers@gmail.com", "Helmi Aziz").
								SetSubject("Email With Attachments").
								SetBody("This is test email with attachments").
								Attach( "cat.webp", catImg ). //<----------------------Add a few pictures as attachments
								AttachFile( "rabbit.jpg", "rabbit.jpg" ).
								AttachFile( "dog.jpg", "dog.jpg" ).
								Dump().
								Send()
    if err != nil {
        log.Fatalf("Failed to send email:\n%v", err)
        return
    }

    log.Println( "Email sent successfully!\nID: ", messageID )
}