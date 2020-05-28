package user

import (
	"github.com/lexbond13/api_core/module/db"
	"github.com/lexbond13/api_core/services/api/handler"
	"github.com/lexbond13/api_core/services/api/handler/validators"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"strconv"
)

type GetHandler struct {
	handler.HandlerBase

	id *int64 `json:"id"`
	response handler.Response
}

// NewGetHandler
func NewGetHandler() handler.IHandler {
	return &GetHandler{}
}

// BindContext
func (g *GetHandler) BindContext(ctx *gin.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return err
	}
	g.id = &id
	return nil
}

// Validate
func (g *GetHandler) Validate() *validators.Validate {
	vErrs := &validators.Validate{}
	vErrs.Required("id", g.id)

	return vErrs
}

// Process
func (g *GetHandler) Process() *handler.Response {

	userDB, err := db.NewUserRepository().FindByID(*g.id)
	if err != nil {
		return g.response.Error(500, err)
	}

	if userDB == nil {
		return g.response.Error(404, errors.New("user not found"))
	}

	return g.response.Success(userDB)
}
