package user

import (
	"github.com/lexbond13/api_core/module/db"
	"github.com/lexbond13/api_core/module/db/structure"
	"github.com/lexbond13/api_core/services/api/handler"
	"github.com/lexbond13/api_core/services/api/handler/validators"
	"github.com/lexbond13/api_core/services/api/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pkg/errors"
)

type AddHandler struct {
	handler.HandlerBase
	request *models.User
	response *handler.Response
}

// NewAddHandler
func NewAddHandler() handler.IHandler {
	return &AddHandler{}
}

// BindContext
func (a *AddHandler) BindContext(ctx *gin.Context) error {
	err := ctx.ShouldBindBodyWith(&a.request, binding.JSON)
	if err != nil {
		return err
	}

	return nil
}

// Validate
func (a *AddHandler) Validate() *validators.Validate {
	return a.request.Validate()
}

// Process
func (a *AddHandler) Process() *handler.Response {
	user := &structure.User{
		Name:   a.request.Name,
		Email:  a.request.Email,
		Status: structure.StatusWaitActivate,
	}

	err := db.NewUserRepository().Create(user)
	if err != nil {
		return a.response.Error(500, errors.Wrap(err, "fail create user"))
	}

	return a.response.Success("user added")
}
