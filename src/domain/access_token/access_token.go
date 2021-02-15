package access_token

import (
	"fmt"
	"github.com/superbkibbles/bookstore_oauth-api/src/utils/errors"
	"github.com/superbkibbles/bookstore_users-api/utils/crypto_utils"
	"strings"
	"time"
)

const (
	expirationTime = 24
	grandTypePassword = "password"
	grandTypeClientCredentials = "client_credentials"
)

type AccessToken struct {
	AccessToken string `json:"access_token"`
	UserID      int64  `json:"user_id"`
	ClientID    int64  `json:"client_id"`
	Expires     int64  `json:"expires"`
}

type AccessTokenRequest struct{
	GrandType string `json:"grand_type"`
	Scope string `json:"scope"`

	// User for password grand type
	Username string `json:"username"`
	Password string `json:"password"`

	// User for client_credentials grand type
	ClientId string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func GetNewAccessToken(userId int64) *AccessToken {
	return &AccessToken{
		UserID: userId,
		Expires: time.Now().UTC().Add(expirationTime * time.Hour).Unix(),
	}
}

func (at *AccessToken) IsExpired() bool {
	return time.Now().UTC().After(time.Unix(at.Expires, 0))
}

func (at *AccessTokenRequest) Validate() *errors.RestErr {
	switch at.GrandType {
	case grandTypePassword:
		break
	case grandTypeClientCredentials:
		break
	default:
		return errors.NewBadRequestErr("invalid grand_type parameter")
	}
	// TODO: Validate parameters for each grand type
	if at.Password == "" {
		return errors.NewBadRequestErr("invalid access token")
	}

	return nil
}

func (at *AccessToken) Generate() {
	at.AccessToken = crypto_utils.GetMd5(fmt.Sprintf("at-%d-%d-ran", at.UserID, at.Expires))
}

func (at *AccessToken) Validate() *errors.RestErr {
	at.AccessToken = strings.TrimSpace(at.AccessToken)
	if at.AccessToken == "" {
		return errors.NewBadRequestErr("invalid access token")
	}
	if at.UserID <= 0 {
		return errors.NewBadRequestErr("invalid User ID")
	}
	if at.ClientID <= 0 {
		return errors.NewBadRequestErr("invalid Client ID")
	}
	if at.Expires <= 0 {
		return errors.NewBadRequestErr("invalid Expiration time")
	}

	return nil
}