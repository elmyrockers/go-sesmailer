# go-sesmailer
[![Go Reference](https://pkg.go.dev/badge/github.com/elmyrockers/go-sesmailer.svg)](https://pkg.go.dev/github.com/elmyrockers/go-sesmailer)
[![Go Version](https://img.shields.io/badge/go1.26+-00ADD8?logo=go&logoColor=white)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Build / CI](https://github.com/elmyrockers/go-sesmailer/actions/workflows/go-ci.yml/badge.svg)](https://github.com/elmyrockers/go-sesmailer/actions/workflows/go-ci.yml)
![Coverage](img/coverage.svg)
![AWS SES](img/amazon-ses.jpg)
<!-- [![Release](https://img.shields.io/github/v/release/elmyrockers/go-sesmailer)](https://github.com/elmyrockers/go-sesmailer/releases) -->

Tiny wrapper around the AWS SES SDK for Go to simplify sending emails with a fluent, chainable API.

***go-sesmailer*** provides a simple and developer-friendly interface to send emails using ***AWS Simple Email Service (SES)***. It supports adding multiple recipients, CC, BCC, Reply-To addresses, and sending both HTML and plain text emails. MIME message construction is handled internally by [go-mimebuilder](https://github.com/elmyrockers/go-mimebuilder), a ***zero-allocation*** in the *`"hot path"`* library, keeping message building ***fast and memory-efficient***. Its fluent API makes it easy to integrate into Go projects, including web frameworks like Fiber.

## Features

- **Lightweight and tiny wrapper** around Amazon SES.
- Built on top of the **official AWS SDK for Go v2**.
- Uses the **Amazon SES API** instead of SMTP for better performance and faster delivery.
- Provides **improved security** by using AWS IAM authentication instead of SMTP credentials.
- **Fluent method chaining** for clean and readable email construction.
- MIME message building powered by ***go-mimebuilder***, a ***secure and high performance*** library.
- Supports **HTML and plain text emails**.
- **Plain text fallback (`AltBody`)** for HTML emails.
- Supports **multiple recipients**: To, CC, and BCC.
- Supports **Reply-To headers**.
- Supports **inline embedded content** (e.g. images) via `Embed`, referenced in HTML using `cid:`.
- Supports **file attachments** via `Attach`.
- Supports **context-based sending** (`SendWithContext`) for cancellation and timeouts.
- Automatically loads **AWS configuration** from the default environment.

## Security Highlights

- **Production-Ready:** Fully tested and safe for use in production environments.
- **Header Injection Protection:** All headers (From, To, Cc, Bcc, Reply-To, Subject) sanitized to strip CR and LF characters, preventing header injection attacks.
- **RFC 2047 Subject Encoding:** Non-ASCII subjects are Q-encoded and safely folded across multiple lines to avoid breaking multi-byte UTF-8 characters.
- **Attachment Security:** Filenames sanitized before use in headers; attachment data base64-encoded with proper MIME line wrapping.
- **Body Encoding:** Quoted-printable encoding (RFC 2045) ensures safe transmission of non-ASCII body content, with line wrapping at 72 characters.
- **AWS SES Secure Delivery:** Uses official SDK v2 with TLS and signed requests.

## Installation
```bash
go get github.com/elmyrockers/go-sesmailer@latest
```
> Make sure you have your AWS credentials set in the environment variables (***AWS_ACCESS_KEY_ID***, ***AWS_SECRET_ACCESS_KEY*** and ***AWS_REGION***) or through your AWS config.
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
import (
    _ "github.com/joho/godotenv/autoload" //Load your .env file automatically
    "github.com/elmyrockers/go-sesmailer"
)

```

## Usage

### 1. Basic Example:

```go
package main

import (
	_ "github.com/joho/godotenv/autoload"
    "github.com/elmyrockers/go-sesmailer"
    "log"
)

func main() {
    builder := sesmailer.New()

    messageID, err := builder.SetFrom("no-reply@yourcompany.com", "Your Company").
                                AddTo("helmi@xeno.com.my", "Helmi Aziz").
                                SetSubject("Test Email").
                                SetBody("Hello! This is a test email.").
                                Send()
    if err != nil {
        log.Fatalf("Failed to send email: %v", err)
        return
    }

    log.Println( "Email sent successfully!\nID: ", messageID )
}

```


### 2. Sending HTML Email with Plain Text Fallback:
```go
builder := sesmailer.New()

messageID, err := builder.SetFrom("no-reply@yourcompany.com", "Your Company").
                            AddTo("helmi@xeno.com.my", "Helmi Aziz").
                            SetSubject("HTML Email Example").
                            SetBody("<h1>Hello</h1><p>This is an HTML email.</p>").
                            AsHTML().
                            SetAltBody("Hello! This is a plain text version.").
                            Send()
if err != nil {
    log.Fatalf("Failed to send email: %v", err)
    return
}
log.Println("Sent! ID:", messageID)
```


### 3. Adding CC, BCC, and Reply-To:
```go
builder := sesmailer.New()

messageID, err := builder.SetFrom("no-reply@yourcompany.com", "Your Company").
                            AddTo("helmi@xeno.com.my", "Helmi Aziz").
                            AddCC("admin@yourcompany.com", "Administrator").
                            AddBCC("your-private-email@gmail.com", "").
                            AddReplyTo("admin@yourcompany.com", "Administrator").
                            SetSubject("Email with CC/BCC/ReplyTo").
                            SetBody("This email has CC, BCC, and Reply-To addresses.").
                        	Send()
if err != nil {
    log.Fatalf("Failed to send email: %v", err)
    return
}
log.Println("Sent! ID:", messageID)
```

### 4. Email with Attachments
```go
package main

import (
    _ "github.com/joho/godotenv/autoload"
    "github.com/elmyrockers/go-sesmailer"
    "log"
    "os"
)

func main() {
    builder := sesmailer.New()

    // Read files as attachments
            invoiceFile, err := os.ReadFile("docs/invoice_123.pdf")
                                    if err != nil {
                                        log.Fatalf("Failed to read attachment: %v", err)
                                        return
                                    }
            logoFile, err := os.ReadFile("images/logo.png")
                                    if err != nil {
                                        log.Fatalf("Failed to read attachment: %v", err)
                                        return
                                    }

    // Build MIME header and body, then send
        messageID, err := builder.SetFrom("no-reply@yourcompany.com", "Your Company").
                                    AddTo("helmi@xeno.com.my", "Helmi Aziz").
                                    SetSubject("Email with Attachments").
                                    SetBody("This email will include a few attachments").
                                    
                                    Attach("invoice.pdf", invoiceFile).
                                    Attach("logo.png", logoFile).
                                    Send()
        if err != nil {
            log.Fatalf("Failed to send email: %v", err)
            return
        }
        log.Println("Sent! ID:", messageID)
}
```

### 5. Dump Raw MIME (For Debugging):
```go
builder := sesmailer.New()

messageID, err := builder.SetFrom("no-reply@yourcompany.com", "Your Company").
                            AddTo("helmi@xeno.com.my", "Helmi Aziz").
                            SetSubject("Debug Email").
                            SetBody("This email will show debug info").
                            Dump().
                            Send()
if err != nil {
    log.Fatalf("Failed to send email: %v", err)
    return
}
log.Println("Sent! ID:", messageID)
```
***
## API Reference

| Method | Description |
|------|-------------|
| `New() *Mailer` | Creates a new `Mailer` instance and automatically initializes the AWS SES client using the default AWS configuration. |
| `SetFrom(email string, name string) *Mailer` | Sets the sender email address and optional display name. |
| `AddTo(email string, name string) *Mailer` | Adds a recipient to the **To** list. |
| `AddAddress(email string, name string) *Mailer` | Alias for `AddTo` - adds a recipient to the **To** list. |
| `AddCC(email string, name string) *Mailer` | Adds a recipient to the **CC** list. |
| `AddBCC(email string, name string) *Mailer` | Adds a recipient to the **BCC** list. |
| `AddReplyTo(email string, name string) *Mailer` | Adds an email address to the **Reply-To** header. |
| `SetSubject(subject string) *Mailer` | Sets the email subject line. |
| `SetBody(body string) *Mailer` | Sets the main email body content. |
| `SetAltBody(alt string) *Mailer` | Sets an alternative plain-text body when sending HTML emails. |
| `AsHTML() *Mailer` | Marks the email content type as `text/html` instead of `text/plain`. |
| `Embed(name string, data []byte, cid string) *Mailer` | Embeds inline content (e.g. an image) referenced in the HTML body via `cid:` - the image displays inline rather than as a downloadable attachment. |
| `Attach(filename string, data []byte) *Mailer` | Adds a file attachment from raw bytes. |
| `Dump() *Mailer` | Builds the MIME message and prints it to stdout, without sending. Useful for debugging. |
| `Send() (string, error)` | Sends the email using a default background context. Returns the SES MessageID and any error. |
| `SendWithContext(ctx context.Context) ( string, error)` | Sends the email using a custom context - useful for timeouts, cancellations, or request tracing. |



## License

This project is licensed under the **MIT License**. See the [LICENSE](LICENSE) file for details.