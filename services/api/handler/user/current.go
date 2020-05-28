package user

import (
	"github.com/lexbond13/api_core/services/api/handler"
	"github.com/lexbond13/api_core/services/api/handler/validators"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type GetCurrentHandler struct {
	handler.HandlerBase
	response *handler.Response
}

func NewGetCurrentHandler() handler.IHandler {
	return &GetCurrentHandler{}
}

func (g *GetCurrentHandler) BindContext(ctx *gin.Context) error {
	g.Ctx = ctx
	return nil
}

func (g *GetCurrentHandler) Validate() *validators.Validate {
	return nil
}

func (g *GetCurrentHandler) Process() *handler.Response {
	identity, err := g.GetIdentity()
	if err != nil {
		return g.response.Error(401, errors.Wrap(err, "fail get identity"))
	}

	return g.response.Success(identity)
}
