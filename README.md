# go-sesmailer

![AWS SES](/amazon-ses.jpg)

Minimal wrapper around the AWS SES SDK for Go to simplify sending emails with a PHPMailer-like API.

go-sesmailer provides a simple and developer-friendly interface to send emails using AWS Simple Email Service (SES). Inspired by PHPMailer, it supports adding multiple recipients, CC, BCC, Reply-To addresses, and sending both HTML and plain text emails. Its fluent API makes it easy to integrate into Go projects, including web frameworks like Fiber.

## Features

- Add multiple recipients, CC, BCC, and Reply-To addresses
- Send plain text or HTML emails
- PHPMailer-style method chaining
- Lightweight and minimal wrapper