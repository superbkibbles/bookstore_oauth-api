package http

import (
	"fmt"
	"net/http"

	access_token2 "github.com/superbkibbles/bookstore_oauth-api/src/services/access_token"
	"github.com/superbkibbles/bookstore_utils-go/rest_errors"

	"github.com/gin-gonic/gin"
	"github.com/superbkibbles/bookstore_oauth-api/src/domain/access_token"
)

type AccessTokenHandler interface {
	GetById(*gin.Context)
	Create(*gin.Context)
	//UpdateExpirationTime(*gin.Context)
}

type accessTokenHandler struct {
	service access_token2.Service
}

func NewHandler(service access_token2.Service) AccessTokenHandler {
	return &accessTokenHandler{
		service: service,
	}
}

func (h *accessTokenHandler) GetById(c *gin.Context) {
	accessToken, err := h.service.GetById(c.Param("access_token_id"))
	if err != nil {
		fmt.Println(err)
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, accessToken)
}

func (h *accessTokenHandler) Create(c *gin.Context) {
	var request access_token.AccessTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		restErr := rest_errors.NewBadRequestErr("Invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	token, err := h.service.Create(request)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusCreated, token)
}

//
//func (h *accessTokenHandler) UpdateExpirationTime(c *gin.Context) {
//
//
//}
