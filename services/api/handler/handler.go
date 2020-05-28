package handler

import (
	"fmt"
	"github.com/lexbond13/api_core/module/logger"
	"github.com/lexbond13/api_core/services/api/handler/validators"
	"github.com/lexbond13/api_core/services/api/models"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

type IHandler interface {
	BindContext(ctx *gin.Context) error
	Validate() *validators.Validate
	Process() *Response
}

type HandlerBase struct {
	Ctx *gin.Context
}

func NewHandler(handlerFunc func() IHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		// bind input params
		response := &Response{}
		handler := handlerFunc()
		err := handler.BindContext(c)
		if err != nil {
			c.JSON(http.StatusOK, response.Error(400, errors.Wrap(err, "fail bind request params")))
			return
		}

		// validators input params
		errs := handler.Validate()
		if errs != nil && errs.Count() > 0 {
			c.JSON(http.StatusOK, response.Error(400, errors.New(fmt.Sprintf("fail validate: %s", errs.String()))))
			return
		}

		// do some work
		resp := handler.Process()
		if !resp.Status {
			logger.Log.Error(resp.err)
		}

		c.JSON(http.StatusOK, resp)
		return
	}
}

func (hb *HandlerBase) GetIdentity() (*models.UserSessionData, error) {
	if hb.Ctx == nil {
		return nil, errors.New("nil context")
	}

	identityKey := "id"
	user, _ := hb.Ctx.Get(identityKey)
	if user == nil {
		return nil, errors.New("no identity found")
	}
	userData := user.(*models.UserSessionData)

	if userData == nil {
		return nil, errors.New("no identity found")
	}
	return userData, nil
}
