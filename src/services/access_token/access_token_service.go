package access_token

import (
	"strings"

	"github.com/superbkibbles/bookstore_oauth-api/src/domain/access_token"
	"github.com/superbkibbles/bookstore_oauth-api/src/repository/db"
	"github.com/superbkibbles/bookstore_oauth-api/src/repository/rest"

	"github.com/superbkibbles/bookstore_utils-go/rest_errors"
)

// Service : acceess token service interface
type Service interface {
	GetById(string) (*access_token.AccessToken, rest_errors.RestErr)
	Create(access_token.AccessTokenRequest) (*access_token.AccessToken, rest_errors.RestErr)
	UpdateExpirationTime(access_token.AccessToken) rest_errors.RestErr
}

type service struct {
	restUsersRepo rest.RestUsersRepository
	dbRepo        db.DbRepository
}

// NewService : creating new service of type Service interface
func NewService(restRepo rest.RestUsersRepository, dbRepo db.DbRepository) Service {
	return &service{
		restUsersRepo: restRepo,
		dbRepo:        dbRepo,
	}
}

// GetById : getting access token by id
func (s *service) GetById(accessTokenID string) (*access_token.AccessToken, rest_errors.RestErr) {
	if len(strings.TrimSpace(accessTokenID)) == 0 {
		return nil, rest_errors.NewBadRequestErr("invalid access token")
	}
	token, err := s.dbRepo.GetById(accessTokenID)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (s *service) Create(request access_token.AccessTokenRequest) (*access_token.AccessToken, rest_errors.RestErr) {
	if err := request.Validate(); err != nil {
		return nil, err
	}
	// TODO: Support both client_credentials and password
	// Authenticate the user against users api
	user, err := s.restUsersRepo.LoginUser(request.Username, request.Password)
	if err != nil {
		return nil, err
	}
	// Generate a new access token
	at := access_token.GetNewAccessToken(user.Id)
	at.Generate()

	// save the new access token in cassandra
	if err := s.dbRepo.Create(*at); err != nil {
		return nil, err
	}
	return at, nil
}

func (s *service) UpdateExpirationTime(at access_token.AccessToken) rest_errors.RestErr {
	if err := at.Validate(); err != nil {
		return err
	}
	return s.dbRepo.UpdateExpirationTime(at)
}
