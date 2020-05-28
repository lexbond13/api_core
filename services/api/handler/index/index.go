package index

import (
	"github.com/lexbond13/api_core/config"
	"github.com/lexbond13/api_core/services/api/handler"
	"github.com/lexbond13/api_core/services/api/handler/validators"
	"github.com/gin-gonic/gin"
)

type MainHandler struct {
	handler.HandlerBase
	response *handler.Response
	params *config.Params
}

// NewMainHandler
func NewMainHandler(params *config.Params) handler.IHandler {
	return &MainHandler{
		params: params,
	}
}

// BindContext
func (m *MainHandler) BindContext(ctx *gin.Context) error {
	return nil
}

// Validate
func (m *MainHandler) Validate() *validators.Validate {
	return &validators.Validate{}
}

// Process
func (m *MainHandler) Process() *handler.Response {
	return m.response.Success(m.params.AppName + " version 1.0")
}
