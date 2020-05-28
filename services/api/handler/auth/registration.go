package auth

import (
	"fmt"
	"github.com/lexbond13/api_core/config"
	"github.com/lexbond13/api_core/module/db"
	"github.com/lexbond13/api_core/module/db/condition"
	"github.com/lexbond13/api_core/module/db/structure"
	"github.com/lexbond13/api_core/services/api/handler"
	"github.com/lexbond13/api_core/services/api/handler/validators"
	"github.com/lexbond13/api_core/services/api/models"
	"github.com/lexbond13/api_core/services/transport/messages"
	"github.com/lexbond13/api_core/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pkg/errors"
)

type AuthHandler struct {
	handler.HandlerBase
	request *models.Auth
	response *handler.Response
	params *config.Params
}

// NewAuthHandler
func NewAuthHandler(config *config.Params) handler.IHandler {
	return &AuthHandler{
		params: config,
	}
}

// BindContext
func (ah *AuthHandler) BindContext(ctx *gin.Context) error {
	err := ctx.ShouldBindBodyWith(&ah.request, binding.JSON)
	if err != nil {
		return err
	}

	return nil
}

// Validate
func (ah *AuthHandler) Validate() *validators.Validate {
	vErrors := &validators.Validate{}

	if !ah.params.OpenAuthMode {
		vErrors.Add(errors.New("open auth mode disabled"))
		return vErrors
	}

	vErrors = ah.request.Validate()

	return vErrors
}

// Process
func (ah *AuthHandler) Process() *handler.Response {

	userRepo := db.NewUserRepository()
	userCondition := condition.NewUserCondition()
	userCondition.SetEmail(ah.request.Email)
	userDB, err := userRepo.FindOne(userCondition)
	if err != nil {
		return ah.response.Error(500, err)
	}

	if userDB == nil || userDB.Status == structure.StatusNew {

		if userDB == nil {
			newUser := structure.User{
				Name:      "user",
				Email:     ah.request.Email,
				ActiveKey: util.RandomString(30),
				Status:    structure.StatusNew,
			}

			err :=  userRepo.Create(&newUser)
			if err != nil {
				return ah.response.Error(500, errors.Wrap(err, "fail create user"))
			}

			userDB = &newUser
		}

		// if message not sent, leave status new for request activate again
		err = ah.sendActivateLink(userDB)
		if err != nil {
			return ah.response.Error(500, errors.Wrap(err, "fail send email message"))
		}

		// update user status for wait activate link
		userDB.Status = structure.StatusWaitActivate
		err = userRepo.Update(userDB)
		if err != nil {
			return ah.response.Error(500, errors.Wrap(err, "fail update user"))
		}

		return ah.response.Success("activation link sent to email")
	}

	if userDB.Status != structure.StatusDisabled {
		// generate new active key for auth another browser
		userDB.ActiveKey = util.RandomString(30)
		err := userRepo.Update(userDB)
		if err != nil {
			return ah.response.Error(500, errors.Wrap(err, "fail update user"))
		}

		err = ah.sendActivateLink(userDB)
		if err != nil {
			return ah.response.Error(500, errors.Wrap(err, "fail send email message"))
		}

		return ah.response.Success("activation link sent to email")
	}

	if userDB.Status == structure.StatusDisabled {
		return ah.response.Error(400, errors.New("user disabled"))
	}

	return ah.response.Success("no actions")
}

// send activate link for user open authorization
func(ah *AuthHandler) sendActivateLink(user *structure.User) error {
	activateURL := fmt.Sprintf("%s/auth/activate?email=%s&active_key=%s", ah.params.ApiURL, user.Email, user.ActiveKey)
	message := messages.NewEmailMessage(messages.GetEmailSender())
	message.ToEmail = user.Email
	message.FromName = ah.params.DomainApp + " Authoriza."

	return message.SendActivateLink(activateURL, ah.params.AppName)
}
