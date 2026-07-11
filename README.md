# go-sesmailer
[![Go Reference](https://pkg.go.dev/badge/github.com/elmyrockers/go-sesmailer.svg)](https://pkg.go.dev/github.com/elmyrockers/go-sesmailer)
[![Go Version](https://img.shields.io/badge/go1.26+-00ADD8?logo=go&logoColor=white)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
![Build / CI](https://github.com/elmyrockers/go-sesmailer/actions/workflows/go-ci.yml/badge.svg)
![Coverage](img/coverage.svg)
![AWS SES](img/amazon-ses.jpg)
<!-- [![Release](https://img.shields.io/github/v/release/elmyrockers/go-sesmailer)](https://github.com/elmyrockers/go-sesmailer/releases) -->

Minimal wrapper around the AWS SES SDK for Go to simplify sending emails with a fluent, chainable API.

**go-sesmailer** provides a simple and developer-friendly interface to send emails using **AWS Simple Email Service (SES)**. It supports adding multiple recipients, CC, BCC, Reply-To addresses, and sending both HTML and plain text emails. MIME message construction is handled internally by [go-mimebuilder](https://github.com/elmyrockers/go-mimebuilder), a **zero-allocation** in the `"hot path"` library, keeping message building ***fast and memory-efficient***. Its fluent API makes it easy to integrate into Go projects, including web frameworks like Fiber.

## Features

- **Lightweight and tiny wrapper** around Amazon SES.
- Built on top of the **official AWS SDK for Go v2**.
- Uses the **Amazon SES API** instead of SMTP for better performance and faster delivery.
<!-- - **High-performance attachments** via **streaming**. Uses `io.Reader` (**32KB chunks**) for **low memory usage**. -->
- Provides **improved security** by using AWS IAM authentication instead of SMTP credentials.
<!-- - **PHPMailer-like API** familiar to developers coming from **PHP and traditional email libraries**. -->
- **Fluent method chaining** for clean and readable email construction.
- Supports **HTML and plain text emails**.
- **Plain text fallback (`AltBody`)** for HTML emails.
- Supports **multiple recipients**: To, CC, and BCC.
- Supports **Reply-To headers**.
<!-- - Built-in **debug logging** with multiple verbosity levels. -->
- Supports **context-based sending** (`SendWithContext`) for cancellation and timeouts.
- Automatically loads **AWS configuration** from the default environment.

## Security Highlights

- **Production-Ready:** Fully tested and safe for use in production environments.
- **Header Injection Protection:** All headers sanitized to remove CR, LF, null, and control characters.
- **RFC 5322 Compliant:** Header lengths truncated at 998 bytes safely, UTF-8 aware
- **Safe Email Addresses:** Validated and properly encoded display names.
- **Attachment Security:** Filenames sanitized; data streamed and base64-encoded.
- **Body Encoding:** Quoted-printable encoding ensures safe transmission of non-ASCII content.
- **AWS SES Secure Delivery:** Uses official SDK v2 with TLS and signed requests.
- **Debug & Logging Control:** Optional debug with minimal exposure of sensitive info.

## Installation
```bash
go get github.com/elmyrockers/go-sesmailer
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
    err := sesmailer.New().
        SetFrom("no-reply@yourcompany.com", "Your Company").
        AddAddress("helmi@xeno.com.my", "Helmi Aziz").
        SetSubject("Test Email").
        SetBody("Hello! This is a test email.").
        IsHTML(false).
        Send()

    if err != nil {
        log.Fatalf("Failed to send email: %v", err)
    }

    log.Println("Email sent successfully")
}

```


### 2. Sending HTML Email with Plain Text Fallback:
```go
err := sesmailer.New().
    SetFrom("no-reply@yourcompany.com", "Your Company").
    AddAddress("helmi@xeno.com.my", "Helmi Aziz").
    SetSubject("HTML Email Example").
    SetBody("<h1>Hello</h1><p>This is an HTML email.</p>").
    SetAltBody("Hello! This is a plain text version.").
    IsHTML(true).
    Send()

if err != nil {
    log.Fatalf("Failed to send email: %v", err)
}
```


### 3. Adding CC, BCC, and Reply-To:
```go
err := sesmailer.New().
    SetFrom("no-reply@yourcompany.com", "Your Company").
    AddAddress("helmi@xeno.com.my", "Helmi Aziz").
    AddCC("admin@yourcompany.com", "Administrator").
    AddBCC("your-private-email@gmail.com", "").
    AddReplyTo("admin@yourcompany.com", "Administrator").
    SetSubject("Email with CC/BCC/ReplyTo").
    SetBody("This email has CC, BCC, and Reply-To addresses.").
	Send()

if err != nil {
    log.Fatalf("Failed to send email: %v", err)
}
```

### 4. Email with Attachments
```go
err := sesmailer.New().
    SetFrom("no-reply@yourcompany.com", "Your Company").
    AddAddress("helmi@xeno.com.my", "Helmi Aziz").
    SetSubject("Email with Attachments").
    SetBody("This email will include a few attachments").
    
    AddAttachment("docs/invoice_123.pdf", "Invoice.pdf").
    AddAttachment("images/logo.png", "CompanyLogo.png").
    Send()

if err != nil {
    log.Fatalf("Failed to send email: %v", err)
}
```

### 5. Enabling Debug Logging:
```go
err := sesmailer.New().
    SetFrom("no-reply@yourcompany.com", "Your Company").
    AddAddress("helmi@xeno.com.my", "Helmi Aziz").
    SetSubject("Debug Email").
    SetBody("This email will show debug info").
    SetDebug(2). // 0 = none, 1 = errors, 2 = verbose
	Send()

if err != nil {
    log.Fatalf("Failed to send email: %v", err)
}
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
| `Send() (*ses.SendRawEmailOutput, error)` | Sends the email using a default background context. Returns the SES API response and any error. |
| `SendWithContext(ctx context.Context) (*ses.SendRawEmailOutput, error)` | Sends the email using a custom context - useful for timeouts, cancellations, or request tracing. |



## License

This project is licensed under the **MIT License**. See the [LICENSE](LICENSE) file for details.