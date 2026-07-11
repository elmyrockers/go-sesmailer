//go:build integration
// +build integration

package sesmailer

import (
	"fmt"
	"time"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"

	_ "github.com/joho/godotenv/autoload"
)

func TestMailer_New(t *testing.T) {
	mailer := New()

	require.NotNil(t, mailer, "New() should never return nil")
	assert.NotNil(t, mailer.client, "expected SES client to be initialized")
	assert.NotNil(t, mailer.builder, "expected mime builder to be initialized")
	assert.Empty(t, mailer.errorList, "expected no errors when AWS config loads successfully")
}

func TestMailer_SendPlainText( t *testing.T ) {
	// Test construct
		mailer := New()
		require.NotNil(t, mailer, "New() should return a non-nil Mailer")

	// Test: add mime headers and body then send
		subject := fmt.Sprintf("sesmailer integration test - plain text - %d", time.Now().UnixNano())
		output, err := mailer.
			SetFrom( "noreply@xeno.com.my", "Xeno System").
			AddReplyTo( "info@xeno.com.my", "Xeno Admin" ).
			AddTo( "elmyrockers@gmail.com", "Developer").
			AddCC( "elmyrockers2@gmail.com", "Maintainer").
			AddBCC( "elmyrockers3@gmail.com", "Project Manager").
			SetSubject(subject).
			
			SetBody("This is a plain text integration test email.").
			Send()

		require.NoError(t, err, "expected Send to succeed")
		require.NotNil(t, output)
		assert.NotEmpty(t, output.MessageId)

		t.Logf("Sent plain text email, MessageId=%s", *output.MessageId)
}

func TestMailer_SendHtmlWithAttachments( t *testing.T ) {
	// Test construct
		mailer := New()
		require.NotNil(t, mailer, "New() should return a non-nil Mailer")

	// Get image data in bytes
		catImg, err := os.ReadFile("/examples/attachment/cat.webp")
			if err != nil { log.Fatal(err) }
		rabbitImg, err := os.ReadFile("/examples/attachment/rabbit.jpg")
			if err != nil { log.Fatal(err) }
		dogImg, err := os.ReadFile("/examples/attachment/dog.jpg")
			if err != nil { log.Fatal(err) }

	// Test: add mime headers and body then send
		subject := fmt.Sprintf("sesmailer integration test - html+attachment - %d", time.Now().UnixNano())
		output, err := mailer.
			SetFrom( "noreply@xeno.com.my", "Xeno System").
			AddReplyTo( "info@xeno.com.my", "Xeno Admin" ).
			AddTo( "elmyrockers@gmail.com", "Developer").
			AddCC( "elmyrockers2@gmail.com", "Maintainer").
			AddBCC( "elmyrockers3@gmail.com", "Project Manager").
			SetSubject(subject).

			SetBody("<h1>Hello</h1><p>This email has an attachment.</p>").AsHTML().
			SetAltBody("Hello. This is the plain text fallback.").
			Attach("cat.webp", catImg).
			Attach("rabbit.jpg", rabbitImg).
			Attach("dog.jpg", dogImg).
			Send()

	require.NoError(t, err, "expected Send to succeed")
	require.NotNil(t, output)
	assert.NotEmpty(t, output.MessageId)

	t.Logf("Sent HTML+attachment email, MessageId=%s", *output.MessageId)
}

func TestMailer_SendHtmlWithEmbeds( t *testing.T ) {

}