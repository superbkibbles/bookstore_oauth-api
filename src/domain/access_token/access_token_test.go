package access_token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetNewAccessTokenConstants(t *testing.T) {
	assert.EqualValues(t, 24, expirationTime, "expiration time should be 24 hours")
}

func TestGetNewAccessToken(t *testing.T) {
	at := GetNewAccessToken(2)
	assert.False(t, at.IsExpired(), "Brand new access token shouldn't be expired")
	assert.EqualValues(t, "", at.AccessToken, "New access token should't have defined access token ID")
	assert.True(t, at.UserID == 2, "New access token should't have an assosiated user id")
}

func TestGetNewAccessTokenIsExpired(t *testing.T) {
	at := AccessToken{}
	assert.True(t, at.IsExpired(), "empty access token should be expired by default")

	at.Expires = time.Now().UTC().Add(3 * time.Hour).Unix()
	assert.False(t, at.IsExpired(), "access token expiring three hours from now should not be expired")
	at.Expires = time.Now().UTC().Add(3 * time.Minute).Unix()
	assert.False(t, at.IsExpired(), "access token expiring three hours from now should not be expired")
}
