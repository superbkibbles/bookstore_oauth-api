package rest

import (
	"fmt"
	"github.com/mercadolibre/golang-restclient/rest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("about to start test cases")
	rest.StartMockupServer()
	os.Exit(m.Run())
}

func TestLoginUserTimeoutFromApi(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "https://api.bookstore.com/users/login",
		ReqBody:     `{"email":"mend@gmail.com","password":"passowrd"}`,
		RespHTTPCode: -1,
		RespBody: `{}`,
	})
	userRep := usersRepository{}
	user, err := userRep.LoginUser("mend@gmail.com", "passowrd")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid restClient response when trying to login user", err.Message)
}

func TestLoginUserInvalidErrorInterface(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "https://api.bookstore.com/users/login",
		ReqBody:     `{"email":"mend@gmail.com","password":"passowrd"}`,
		RespHTTPCode: http.StatusNotFound,
		RespBody: `{"message": "invalid login credentials", "status": "404", "error": "not_found"}`,
	})
	userRep := usersRepository{}
	user, err := userRep.LoginUser("mend@gmail.com", "passowrd")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "Invalid error interface when trying to login user", err.Message)
}

func TestLoginUserInvalidLoginCredentials(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "https://api.bookstore.com/users/login",
		ReqBody:     `{"email":"mend@gmail.com","password":"passowrd"}`,
		RespHTTPCode: http.StatusNotFound,
		RespBody: `{"message": "invalid login credentials", "status": 404, "error": "not_found"}`,
	})
	userRep := usersRepository{}
	user, err := userRep.LoginUser("mend@gmail.com", "passowrd")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.Status)
	assert.EqualValues(t, "invalid login credentials", err.Message)
}

func TestLoginUserInvalidUserJsonResponse(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "https://api.bookstore.com/users/login",
		ReqBody:     `{"email":"mend@gmail.com","password":"passowrd"}`,
		RespHTTPCode: http.StatusOK,
		RespBody: `{"id":"6","first_name":"Mand","last_name":"haj ali","email":"mand2@gmail.com","date_created":"2021-02-04 11:11:41","status":"active"}`,
	})
	userRep := usersRepository{}
	user, err := userRep.LoginUser("mend@gmail.com", "passowrd")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "error while trying to unmarshal users login response", err.Message)
}

func TestLoginUserNoError(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "https://localhost:8080/users/login",
		ReqBody:     `{"email":"mend@gmail.com","password":"passowrd"}`,
		RespHTTPCode: http.StatusOK,
		RespBody: `{"id":6,"first_name":"Mand","last_name":"haj ali","email":"mand2@gmail.com","date_created":"2021-02-04 11:11:41","status":"active"}`,
	})
	userRep := usersRepository{}
	user, err := userRep.LoginUser("mend@gmail.com", "passowrd")
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.EqualValues(t, 6, user.Id)
	assert.EqualValues(t, "Mand", user.FirstName)
	assert.EqualValues(t, "haj ali", user.LastName)
	assert.EqualValues(t, "mand2@gmail.com", user.Email)
}