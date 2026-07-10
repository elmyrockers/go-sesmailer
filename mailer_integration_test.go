//go:build integration
// +build integration

package sesmailer

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
)

func TestMailer_New(t *testing.T) {
	mailer := New()

	require.NotNil(t, mailer, "New() should never return nil")
	assert.NotNil(t, mailer.client, "expected SES client to be initialized")
	assert.NotNil(t, mailer.builder, "expected mime builder to be initialized")
	assert.Empty(t, mailer.errorList, "expected no errors when AWS config loads successfully")
}

func TestMailer_SendPlainText( t *testing.T ) {

}

func TestMailer_SendHtmlWithAttachments( t *testing.T ) {
	
}

func TestMailer_SendHtmlWithEmbeds( t *testing.T ) {

}