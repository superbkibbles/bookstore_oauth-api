package access_token

import (
	"github.com/superbkibbles/bookstore_oauth-api/src/domain/access_token"
	"github.com/superbkibbles/bookstore_oauth-api/src/repository/rest"
	"strings"

	"github.com/superbkibbles/bookstore_oauth-api/src/utils/errors"
)

type Service interface {
	GetById(string) (*access_token.AccessToken, *errors.RestErr)
	Create(access_token.AccessTokenRequest) (*access_token.AccessToken, *errors.RestErr)
	UpdateExpirationTime(access_token.AccessToken) *errors.RestErr
}

type Repository interface {
	GetById(string) (*access_token.AccessToken, *errors.RestErr)
	Create(token access_token.AccessToken) *errors.RestErr
	UpdateExpirationTime(access_token.AccessToken) *errors.RestErr
}

type service struct {
	restUsersRepo rest.RestUsersRepository
	dbRepo Repository
}

func NewService(restRepo rest.RestUsersRepository, dbRepo Repository) Service {
	return &service{
		restUsersRepo: restRepo,
		dbRepo: dbRepo,
	}
}

func (s *service) GetById(accessTokenId string) (*access_token.AccessToken, *errors.RestErr) {
	if len(strings.TrimSpace(accessTokenId)) == 0 {
		return nil, errors.NewBadRequestErr("invalid access token")
	}
	token, err := s.dbRepo.GetById(accessTokenId)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (s *service) Create(request access_token.AccessTokenRequest) (*access_token.AccessToken, *errors.RestErr) {
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
func (s *service) UpdateExpirationTime(at access_token.AccessToken) *errors.RestErr {
	if err := at.Validate(); err != nil {
		return err
	}
	return s.dbRepo.UpdateExpirationTime(at)
}