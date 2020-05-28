package server

import (
	"fmt"
	"github.com/lexbond13/api_core/config"
	logger2 "github.com/lexbond13/api_core/module/logger"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Server struct {
	Config *config.Server
	Params *config.Params
	logger logger2.ILogger
}

// NewServer
func NewServer(config *config.Server, params *config.Params, logger logger2.ILogger) *Server {
	return &Server{
		Config: config,
		Params: params,
		logger: logger,
	}
}

// Run
func (s *Server) Run() error {

	// Set GIN mode
	if s.Params.AppDevMode {
		s.logger.Info("WARNING: starting server on DEV MODE")
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// load all routes and depend handlers
	router, err := GetRouter(s.Params)
	if err != nil {
		return errors.Wrap(err, "fail get router")
	}

	addr := fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port)

	s.logger.Info(fmt.Sprintf("Starting server(%s) on address %s", s.Params.NodeName, addr))
	if err := router.Run(addr); err != nil {
		return errors.Wrap(err, "starting server failed")
	}

	return nil
}
