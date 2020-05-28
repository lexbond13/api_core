package server

import (
	"github.com/lexbond13/api_core/config"
	"github.com/lexbond13/api_core/services/api/handler"
	"github.com/lexbond13/api_core/services/api/handler/auth"
	"github.com/lexbond13/api_core/services/api/handler/club"
	"github.com/lexbond13/api_core/services/api/handler/file"
	"github.com/lexbond13/api_core/services/api/handler/index"
	"github.com/lexbond13/api_core/services/api/handler/user"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// GetRouter
func GetRouter(params *config.Params) (*gin.Engine, error) {

	router := gin.Default()

	// depend middlewares
	SetUpMiddlewares(router, params)
	authMiddleware, err := NewAuthMiddleware(params)
	if err != nil {
		return nil, errors.Wrap(err, "JWT Error")
	}

	// Default endpoint
	router.GET("/", handler.NewHandler(NewParamsWrap(params, index.NewMainHandler)))

	// Clubs
	clubGroup := router.Group("/club")
	{
		clubGroup.Use(authMiddleware.MiddlewareFunc())
		clubGroup.POST("/", handler.NewHandler(NewParamsWrap(params, club.NewAddHandler)))
		clubGroup.DELETE("/:id", handler.NewHandler(club.NewDeleteHandler))
	}

	// Authenticate
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/registration", handler.NewHandler(NewParamsWrap(params, auth.NewAuthHandler)))
		authGroup.GET("/activate", handler.NewHandler(NewParamsWrap(params, auth.NewAuthActivateHandler)))

		authGroup.Use(authMiddleware.MiddlewareFunc())
		authGroup.GET("/refresh_token", handler.NewHandler(NewParamsWrap(params, auth.NewRefreshTokenHandler)))
	}

	// Users
	userGroup := router.Group("/user")
	userGroup.Use(authMiddleware.MiddlewareFunc())
	{
		userGroup.GET("/", handler.NewHandler(user.NewGetCurrentHandler))
	}

	// Files
	fileGroup := router.Group("/files")
	fileGroup.Use(authMiddleware.MiddlewareFunc())
	{
		fileGroup.POST("/upload", handler.NewHandler(NewParamsWrap(params, file.NewUploadHandler)))
	}

	// Handle 404 route
	router.NoRoute(func(c *gin.Context) {
		response := handler.Response{}
		c.JSON(200, response.Error(404, errors.New("page not found")))
	})

	return router, nil
}

func NewParamsWrap(params *config.Params, hand func(params *config.Params)  handler.IHandler) func() handler.IHandler {
	return func() handler.IHandler {
		fn := hand(params)
		return fn
	}
}
