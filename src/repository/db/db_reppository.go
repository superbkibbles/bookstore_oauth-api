package db

import (
	"github.com/gocql/gocql"
	"github.com/superbkibbles/bookstore_oauth-api/src/client/cassandra"
	"github.com/superbkibbles/bookstore_oauth-api/src/domain/access_token"
	"github.com/superbkibbles/bookstore_oauth-api/src/utils/errors"
)

const (
	queryGetAccessToken = "SELECT access_token, user_id, client_id, expires from access_tokens WHERE access_token=?;"
	queryCreateAccessToken = "INSERT into access_tokens(access_token, user_id, client_id, expires) values(?,?,?,?);"
	queryUpdateExpires = "UPDATE access_tokens SET expires=? WHERE access_token=?;"
)

// DbRepository is interface of db_repository
type DbRepository interface {
	GetById(string) (*access_token.AccessToken, *errors.RestErr)
	Create(access_token.AccessToken) *errors.RestErr
	UpdateExpirationTime(access_token.AccessToken) *errors.RestErr
}

type dbRepository struct{}

func NewRepository() DbRepository {
	return &dbRepository{}
}

func (r *dbRepository) GetById(id string) (*access_token.AccessToken, *errors.RestErr) {
	var result access_token.AccessToken
	if err := cassandra.GetSession().Query(queryGetAccessToken, id).Scan(
		&result.AccessToken,
		&result.UserID,
		&result.ClientID,
		&result.Expires,
		); err != nil {
		if err == gocql.ErrNotFound {
			return nil, errors.NewNotFoundErr("no access token found with given id")
		}
		return nil, errors.NewInternalServerErr(err.Error())
	}
	return &result, nil
}

func (r *dbRepository) Create(at access_token.AccessToken) *errors.RestErr {
	if err  :=cassandra.GetSession().Query(queryCreateAccessToken,
		at.AccessToken,
		at.UserID,
		at.ClientID,
		at.Expires,
	).Exec(); err != nil {
		return errors.NewInternalServerErr(err.Error())
	}
	return nil
}

func (r *dbRepository) UpdateExpirationTime(at access_token.AccessToken) *errors.RestErr {
	if err  :=cassandra.GetSession().Query(queryUpdateExpires,
		at.Expires,
		at.AccessToken,
	).Exec(); err != nil {
		return errors.NewInternalServerErr(err.Error())
	}
	return nil
}