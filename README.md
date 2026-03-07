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
> Make sure you have your AWS credentials set in the environment (***AWS_ACCESS_KEY_ID*** and ***AWS_SECRET_ACCESS_KEY***) or through your AWS config.

## Usage

### Basic Example

```go
package main

import (
    "log"
    "github.com/elmyrockers/go-sesmailer"
)

func main() {
    mail := sesmailer.New().
        SetFrom("no-reply@yourcompany.com", "Your Company").
        AddAddress("helmi@xeno.com.my", "Helmi Aziz").
        SetSubject("Test Email").
        SetBody("Hello! This is a test email.").
        IsHTML(false)

    if err := mail.Send(); err != nil {
        log.Fatalf("Failed to send email: %v", err)
    }

    log.Println("Email sent successfully")
}

```

***
### Sending HTML Email with Plain Text Fallback
```go
mail := sesmailer.New().
    SetFrom("no-reply@yourcompany.com", "Your Company").
    AddAddress("helmi@xeno.com.my", "Helmi Aziz").
    SetSubject("HTML Email Example").
    SetBody("<h1>Hello</h1><p>This is an HTML email.</p>").
    SetAltBody("Hello! This is a plain text version.").
    IsHTML(true)

if err := mail.Send(); err != nil {
    log.Fatalf("Failed to send email: %v", err)
}
```

***
### Adding CC, BCC, and Reply-To
```go
mail := sesmailer.New().
    SetFrom("no-reply@yourcompany.com", "Your Company").
    AddAddress("helmi@xeno.com.my", "Helmi Aziz").
    AddCC("admin@yourcompany.com", "Administrator").
    AddBCC("your-private-email@gmail.com", "").
    AddReplyTo("admin@yourcompany.com", "Administrator").
    SetSubject("Email with CC/BCC/ReplyTo").
    SetBody("This email has CC, BCC, and Reply-To addresses.").
	Send()
```

***
### Enabling Debug Logging
```go
mail := sesmailer.New().
    SetFrom("no-reply@yourcompany.com", "Your Company").
    AddAddress("helmi@xeno.com.my", "Helmi Aziz").
    SetSubject("Debug Email").
    SetBody("This email will show debug info").
    SetDebug(2). // 0 = none, 1 = errors, 2 = verbose
	Send()
```