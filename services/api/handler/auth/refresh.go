package auth

import (
	"fmt"
	"github.com/lexbond13/api_core/config"
	"github.com/lexbond13/api_core/module/cache"
	"github.com/lexbond13/api_core/services/api/handler"
	"github.com/lexbond13/api_core/services/api/handler/validators"
	"github.com/lexbond13/api_core/services/api/models"
	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type RefreshTokenHandler struct {
	handler.HandlerBase
	response *handler.Response
	params *config.Params
}

func NewRefreshTokenHandler(params *config.Params) handler.IHandler {
	return &RefreshTokenHandler{
		params: params,
	}
}

func (r *RefreshTokenHandler) BindContext(ctx *gin.Context) error {
	r.Ctx = ctx
	return nil
}

func (r *RefreshTokenHandler) Validate() *validators.Validate {
	//claims, err := mw.CheckIfTokenExpire(c)
	//if err != nil {
	//	return "", time.Now(), err
	//}
	return nil
}

func (r *RefreshTokenHandler) Process() *handler.Response {

	token, exists := r.Ctx.Get("JWT_TOKEN")
	if !exists {
		return r.response.Error(500, errors.New("token not exist"))
	}

	oldToken := token.(string)

	claimsModel, err := ParseToken(oldToken, r.params.AuthSecret)
	if err != nil {
		return r.response.Error(500, errors.Wrap(err, "fail parse token"))
	}

	expire := jwtGo.TimeFunc().Add(time.Duration(r.params.MaxCookieLifeTimeHours) * time.Hour)
	claimsModel.ExpiresAt = expire.Unix()

	tokenClaims := jwtGo.NewWithClaims(jwtGo.SigningMethodHS256, claimsModel)

	newToken, err := tokenClaims.SignedString([]byte(r.params.AuthSecret))
	if err != nil {
		return r.response.Error(500, errors.Wrap(err, "fail refresh token"))
	}

	appName := strings.ReplaceAll(r.params.AppName, " ", "")
	maxage := int(expire.Unix() - time.Now().Unix())

	// if cache storage exist, update key there
	if cacheClient := cache.GetClient(); cacheClient != nil {
		userSessionData, err := cacheClient.Get(oldToken)
		if err != nil {
			return r.response.Error(500, errors.Wrap(err, "fail get token from storage"))
		}

		err = cacheClient.Set(newToken, userSessionData, maxage)
		if err != nil {
			return r.response.Error(500, errors.Wrap(err, "fail set token from storage"))
		}

		ok, err := cacheClient.Del(oldToken)
		if !ok || err != nil {
			// expired self or add this to log
		}
	}

	// set cookie
	r.Ctx.SetCookie(
		appName,
		newToken,
		maxage,
		"/",
		r.params.DomainApp,
		false,
		false,
	)

	return r.response.Success(fmt.Sprintf("token refreshed. expire: %s", expire))
}

func ParseToken(token, authSecret string) (*models.Claims, error) {
	tokenClaims, err := jwtGo.ParseWithClaims(token, &models.Claims{}, func(token *jwtGo.Token) (interface{}, error) {
		return []byte(authSecret), nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*models.Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
