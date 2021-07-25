package rest

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/mercadolibre/golang-restclient/rest"
	"github.com/superbkibbles/bookstore_oauth-api/src/domain/users"
	"github.com/superbkibbles/bookstore_utils-go/rest_errors"
)

var (
	usersRestClient = rest.RequestBuilder{
		BaseURL: "http://localhost:8080",
		Timeout: 100 * time.Millisecond,
	}
)

type RestUsersRepository interface {
	LoginUser(string, string) (*users.User, rest_errors.RestErr)
}

type usersRepository struct{}

func NewRepository() RestUsersRepository {
	return &usersRepository{}
}

func (u *usersRepository) LoginUser(email, password string) (*users.User, rest_errors.RestErr) {
	request := users.UserLoginRequest{
		Email:    email,
		Password: password,
	}

	response := usersRestClient.Post("/users/login", request)
	if response == nil || response.Response == nil {
		return nil, rest_errors.NewInternalServerErr("invalid restClient response when trying to login user", errors.New("restClient error"))
	}
	if response.StatusCode > 299 {
		apiErr, err := rest_errors.NewRestErrorFromBytes(response.Bytes())
		if err != nil {
			return nil, rest_errors.NewInternalServerErr("Invalid error interface when trying to login user", errors.New("ds"))
		}
		return nil, apiErr
	}
	var user users.User
	if err := json.Unmarshal(response.Bytes(), &user); err != nil {
		return nil, rest_errors.NewInternalServerErr("error while trying to unmarshal users login response", errors.New("JSON parsing error"))
	}
	return &user, nil
}
