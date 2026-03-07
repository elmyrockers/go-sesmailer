# go-sesmailer

![AWS SES](/amazon-ses.jpg)

Minimal wrapper around the AWS SES SDK for Go to simplify sending emails with a PHPMailer-like API.

go-sesmailer provides a simple and developer-friendly interface to send emails using AWS Simple Email Service (SES). Inspired by PHPMailer, it supports adding multiple recipients, CC, BCC, Reply-To addresses, and sending both HTML and plain text emails. Its fluent API makes it easy to integrate into Go projects, including web frameworks like Fiber.

## Features

- Add multiple recipients, CC, BCC, and Reply-To addresses
- Send plain text or HTML emails
- PHPMailer-style method chaining
- Lightweight and minimal wrapper








## Installation
```bash
go get github.com/elmyrockers/go-sesmailer
```
Make sure you have your AWS credentials set in the environment (AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY) or through your AWS config.

## Usage

### Basic Example

```go
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

```