package e2e

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2/clientcredentials"
)

// The CLI logs in with golang.org/x/oauth2, which form-encodes the secret in the Basic header
// (RFC 6749 §2.3.1). createTestUser's password contains "!", which it sends as "%21".
func TestOAuthClientCredentials_SecretWithSpecialChars(t *testing.T) {
	env := setupTestEnv(t)
	defer teardownTestEnv(t, env)

	u := createTestUser(t, env)
	cfg := clientcredentials.Config{
		ClientID:     u.Username,
		ClientSecret: u.Password,
		TokenURL:     env.BaseURL + "/v1/oauth/token",
		Scopes:       []string{"socialapp.users.read"},
	}
	tok, err := cfg.Token(t.Context())
	require.NoError(t, err)
	require.NotEmpty(t, tok.AccessToken)

	cfg.ClientSecret = u.Password + "x"
	_, err = cfg.Token(t.Context())
	require.Error(t, err, "a wrong secret must not get a token")
}
