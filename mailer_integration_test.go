// +build integration

package sesmailer

import (
	// "os"
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	_ "github.com/joho/godotenv/autoload"
	// "github.com/davecgh/go-spew/spew"
)

func TestIntegration_New(t *testing.T) {
	// Skipping test in short mode
	    if testing.Short() {
	        t.Skip("Skipping this test in short mode")
	    }

	// Load AWS config manually, like New() would
			ctx := context.Background()
			cfg, err := config.LoadDefaultConfig(ctx)
			if err != nil {
				t.Fatalf("Unable to load AWS config: %v", err)
			}

			creds, err := cfg.Credentials.Retrieve(ctx)
			if err != nil {
				t.Fatalf("Unable to retrieve AWS credentials: %v", err)
			}

			if creds.AccessKeyID == "" || creds.SecretAccessKey == "" {
				t.Fatal("AWS credentials are empty")
			}

			if cfg.Region == "" {
				t.Fatal("AWS region is not set")
			}


	// Now create Mailer and check defaults
			m := New()
			if m == nil {
				t.Fatal("expected Mailer instance, got nil")
			}

			if m.client == nil {// Make sure the SES client is initialized
				t.Fatal("expected SES client to be initialized")
			}

			// Check default values
			if m.ContentType != "text/plain" {
				t.Errorf("expected default ContentType text/plain, got %s", m.ContentType)
			}
			if m.Debug != 0 {
				t.Errorf("expected default Debug = 0, got %d", m.Debug)
			}

			// Check that all slices are initialized and empty
			if len(m.To) != 0 || len(m.Cc) != 0 || len(m.Bcc) != 0 || len(m.ReplyTo) != 0 {
				t.Errorf("expected empty slices; got To:%v Cc:%v Bcc:%v ReplyTo:%v", m.To, m.Cc, m.Bcc, m.ReplyTo)
			}
}


// TestIntegration_SendContext tests sending email with SendContext
func TestIntegration_SendContext(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	mailer := New().
		SetFrom("no-reply@xeno.com.my", "Xeno System").
		AddAddress("elmyrockers@gmail.com", "Helmi Aziz").
		SetSubject("Integration Test SendContext()").
		SetBody("This is a test email from SendContext() integration test").
		IsHTML(false).
		SetDebug(2)

	ctx := context.Background()
	err := mailer.SendContext(ctx)
	if err != nil {
		t.Fatalf("SendContext() failed: %v", err)
	}
}

// TestIntegration_Send tests sending email with Send (default background context)
func TestIntegration_Send(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	err := New().
		SetFrom("no-reply@xeno.com.my", "Xeno System").
		AddAddress("elmyrockers@gmail.com", "Helmi Aziz").
		SetSubject("Integration Test Send()").
		SetBody("<p>This is a test email from Send() integration test</p>").
		SetAltBody("This is a test email: Alt body").
		IsHTML(true).
		SetDebug(2).
		Send()

	if err != nil {
		t.Fatalf("Send() failed: %v", err)
	}
}

func TestIntegration_Send_NoRecipient(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

    err := New().
		SetFrom("no-reply@xeno.com.my", "").
		SetSubject("Test").
		SetBody("Test body").
		Send()

    if err == nil {
        t.Fatal("expected error but got nil")
    }
}