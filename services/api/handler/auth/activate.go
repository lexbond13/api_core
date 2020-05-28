package auth

import (
	"encoding/json"
	"github.com/lexbond13/api_core/config"
	"github.com/lexbond13/api_core/module/cache"
	"github.com/lexbond13/api_core/module/db"
	"github.com/lexbond13/api_core/module/db/condition"
	"github.com/lexbond13/api_core/module/db/structure"
	"github.com/lexbond13/api_core/services/api/handler"
	"github.com/lexbond13/api_core/services/api/handler/validators"
	"github.com/lexbond13/api_core/services/api/models"
	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strings"
	"time"
)

type AuthActivateHandler struct {
	handler.HandlerBase
	request *models.AuthActivate
	response *handler.Response
	params *config.Params
}

func NewAuthActivateHandler(params *config.Params) handler.IHandler {
	return &AuthActivateHandler{
		params: params,
	}
}

func (a *AuthActivateHandler) BindContext(ctx *gin.Context) error {
	a.HandlerBase.Ctx = ctx
	a.request = &models.AuthActivate{
		Email:     ctx.Query("email"),
		ActiveKey: ctx.Query("active_key"),
	}

	return nil
}

func (a *AuthActivateHandler) Validate() *validators.Validate {
	return a.request.Validate()
}

func (a *AuthActivateHandler) Process() *handler.Response {

	userRepo := db.NewUserRepository()
	userCondition := condition.NewUserCondition()
	userCondition.SetEmail(a.request.Email)
	userDB, err := userRepo.FindOne(userCondition)
	if err != nil {
		return a.response.Error(500, errors.Wrap(err, "fail find user"))
	}

	if userDB == nil {
		return a.response.Error(400, errors.New("user not found"))
	}

	if userDB.Status == structure.StatusDisabled {
		return a.response.Error(400, errors.New("user disabled"))
	}

	if userDB.ActiveKey != a.request.ActiveKey {
		return a.response.Error(400, errors.New("active key is wrong"))
	}

	userDB.ActiveKey = ""
	userDB.Status = structure.StatusEnabled
	err = userRepo.Update(userDB)
	if err != nil {
		return a.response.Error(500, errors.Wrap(err, "fail update user"))
	}

	appName := strings.ReplaceAll(a.params.AppName, " ", "")
	expire := jwtGo.TimeFunc().Add(time.Duration(a.params.MaxCookieLifeTimeHours) * time.Hour)
	claims := &models.Claims{
		ID:    userDB.ID,
		Email: userDB.Email,
		StandardClaims: jwtGo.StandardClaims {
			ExpiresAt : expire.Unix(),
			Issuer : appName,
		},
	}

	tokenClaims := jwtGo.NewWithClaims(jwtGo.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(a.params.AuthSecret))

	if err != nil {
		return a.response.Error(500, errors.Wrap(err, "fail get token for user"))
	}

	// set cookie here
	maxage := int(expire.Unix() - time.Now().Unix())
	a.Ctx.SetCookie(
		appName,
		token,
		maxage,
		"/",
		a.params.DomainApp,
		false,
		false,
	)

	// set identity to cache
	userSession, err := json.Marshal(&models.UserSessionData{
		ID:    userDB.ID,
		Email: userDB.Email,
		Status: userDB.Status,
	})

	if err != nil {
		return a.response.Error(500, errors.Wrap(err, "fail set data to cache"))
	}

	if cacheClient := cache.GetClient(); cacheClient != nil {
		err = cache.GetClient().Set(token, string(userSession), maxage)
	}

	if err != nil {
		return a.response.Error(500, errors.Wrap(err, "fail set data to cache"))
	}

	a.HandlerBase.Ctx.Redirect(http.StatusFound, a.params.FrontendURL)
	return a.response.Success("Authorized")
}
