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