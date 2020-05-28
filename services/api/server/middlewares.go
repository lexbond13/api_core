package server

import (
	"encoding/json"
	"github.com/lexbond13/api_core/config"
	"github.com/lexbond13/api_core/module/cache"
	"github.com/lexbond13/api_core/module/db/structure"
	"github.com/lexbond13/api_core/services/api/models"
	jwt "github.com/appleboy/gin-jwt/v2"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

var identityKey = "id"

// SetUpMiddlewares
func SetUpMiddlewares(router *gin.Engine, params *config.Params) {
	// CORS
	if !params.AppDevMode {
		configCors := cors.DefaultConfig()
		configCors.AllowOrigins = []string{params.FrontendURL, params.ApiURL}
		configCors.AllowMethods = []string{"OPTIONS", "PUT", "DELETE"}
		configCors.AllowCredentials = true
		router.Use(cors.New(configCors))
	}

	// Sentry depend
	router.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))

	router.Use(func(ctx *gin.Context) {
		if hub := sentrygin.GetHubFromContext(ctx); hub != nil {
			hub.Scope().SetTag(params.AppName, params.DomainApp)
		}
		ctx.Next()
	})
}

// NewAuthMiddleware authenticator, authorizator and login
func NewAuthMiddleware(params *config.Params) (authMiddleware *jwt.GinJWTMiddleware, err error) {
	// the jwt middleware
	appName := strings.ReplaceAll(params.AppName, " ", "")
	authMiddleware, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:       params.AppName,
		Key:         []byte(params.AuthSecret),
		Timeout:     time.Duration(params.MaxCookieLifeTimeHours) * time.Hour, // set max token live time this
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			return nil
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
					return &models.UserSessionData{
						ID: int64(claims["id"].(float64)),
						Email: claims["email"].(string),
			}
		},
		SendCookie:     true,
		SecureCookie:   false, //non HTTPS dev environments
		CookieHTTPOnly: false,  // JS can't modify
		CookieDomain:   params.DomainApp,
		CookieName:     appName,
		TokenLookup:    "cookie:"+appName,
		Authenticator: func(c *gin.Context) (interface{}, error) {
			return nil, nil
		},

		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*models.UserSessionData); ok {
				// Check user data in storage for allow or denied access here (example: user status) if storage is set
				if cacheClient := cache.GetClient(); cacheClient != nil {
					if token, ok := c.Get("JWT_TOKEN"); ok {
						userCache, err := cache.GetClient().Get(token.(string))
						if err != nil {
							return false
						}
						userSession := &models.UserSessionData{}
						err = json.Unmarshal([]byte(userCache), &userSession)
						if err != nil {
							return false
						}
						// add conditions for check user data from token and cache
						if userSession.Status == structure.StatusDisabled {
							return false
						}
						return true
					}
				}

				return true
			}
			return false
		},

		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(http.StatusOK, gin.H{
				"Data":    "401 Unauthorized or user undefined",
				"Status": false,
				"Errors": []int64{401},
			})
		},
		LoginResponse: func(c *gin.Context, code int, token string, t time.Time) {
			return
		},
		RefreshResponse:  func(c *gin.Context, code int, token string, t time.Time) {
			c.JSON(http.StatusOK, gin.H{
				"Data":    "login successfully",
				"Status": true,
				"Errors": []int64{0},
			})
		},
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	return
}
