# go-sesmailer

![AWS SES](/amazon-ses.jpg)

Minimal wrapper around the AWS SES SDK for Go to simplify sending emails with a PHPMailer-like API.

go-sesmailer provides a simple and developer-friendly interface to send emails using AWS Simple Email Service (SES). Inspired by PHPMailer, it supports adding multiple recipients, CC, BCC, Reply-To addresses, and sending both HTML and plain text emails. Its fluent API makes it easy to integrate into Go projects, including web frameworks like Fiber.

## Features

- **Lightweight and tiny wrapper** around Amazon SES.
- Built on top of the **official AWS SDK for Go v2**.
- **PHPMailer-like API** familiar to developers coming from **PHP and traditional email libraries**.
- **Fluent method chaining** for clean and readable email construction.
- Supports **HTML and plain text emails**.
- **Plain text fallback (`AltBody`)** for HTML emails.
- Supports **multiple recipients**: To, CC, and BCC.
- Supports **Reply-To headers**.
- Built-in **debug logging** with multiple verbosity levels.
- Supports **context-based sending** (`SendContext`) for cancellation and timeouts.
- Automatically loads **AWS configuration** from the default environment.


## Installation
```bash
go get github.com/elmyrockers/go-sesmailer
```
> Make sure you have your AWS credentials set in the environment variables (***AWS_ACCESS_KEY_ID*** and ***AWS_SECRET_ACCESS_KEY***) or through your AWS config.
> During development, you can create a `.env` file next to your `main.go` and load it using the `joho/godotenv` library:
```bash
go get github.com/joho/godotenv
```
> Your `.env` file should contain:
```env
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
AWS_REGION=us-east-1
```
> Then you can load it automatically in your code using the import:
```go
package main

import (
    _ "github.com/joho/godotenv/autoload"
    "github.com/elmyrockers/go-sesmailer"
)

```

## Usage

### 1. Basic Example

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


### 2. Sending HTML Email with Plain Text Fallback
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


### 3. Adding CC, BCC, and Reply-To
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


### 4. Enabling Debug Logging
```go
mail := sesmailer.New().
    SetFrom("no-reply@yourcompany.com", "Your Company").
    AddAddress("helmi@xeno.com.my", "Helmi Aziz").
    SetSubject("Debug Email").
    SetBody("This email will show debug info").
    SetDebug(2). // 0 = none, 1 = errors, 2 = verbose
	Send()
```
***
## API Reference

| Method | Description |
|------|-------------|
| `New() *Mail` | Creates a new `Mail` instance and automatically initializes the AWS SES client using the default AWS configuration. |
| `SetFrom(email string, name string) *Mail` | Sets the sender email address and optional display name. |
| `AddAddress(email string, name string) *Mail` | Adds a recipient to the **To** list. |
| `AddCC(email string, name string) *Mail` | Adds a recipient to the **CC** list. |
| `AddBCC(email string, name string) *Mail` | Adds a recipient to the **BCC** list. |
| `AddReplyTo(email string, name string) *Mail` | Adds an email address to the **Reply-To** header. |
| `SetSubject(subject string) *Mail` | Sets the email subject line. |
| `SetBody(body string) *Mail` | Sets the main email body content. |
| `SetAltBody(alt string) *Mail` | Sets an alternative plain-text body when sending HTML emails. |
| `IsHTML(isHtml bool) *Mail` | Sets whether the email content type should be `text/html` or `text/plain`. |
| `SetDebug(level int) *Mail` | Enables debug logging. `0 = disabled`, `1 = errors/retries`, `2 = verbose request/response logs`. |
| `Send() error` | Sends the email using a default background context. |
| `SendContext(ctx context.Context) error` | Sends the email using a custom context. Useful for timeouts, cancellations, or request tracing. |