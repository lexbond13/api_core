package main

import (
	"github.com/lexbond13/api_core/config"
	"github.com/lexbond13/api_core/module/cache"
	"github.com/lexbond13/api_core/module/db"
	"github.com/lexbond13/api_core/module/logger"
	"github.com/lexbond13/api_core/services/api/server"
	"github.com/lexbond13/api_core/services/transport/files"
	"github.com/lexbond13/api_core/services/transport/messages"
	"github.com/pkg/errors"
)

func main() {

	appConfig, err := config.Init()
	if err != nil {
		panic(err)
	}

	err = Init(appConfig)
	if err != nil {
		panic(err)
	}

	webServer := server.NewServer(appConfig.Server, appConfig.Params, logger.Log)

	err = webServer.Run()
	if err != nil {
		logger.Log.Fatal(err)
	}
}

// Init init all modules
func Init(config *config.AppConfig) error {

	err := logger.Init(config.Logger, config.Params.AppDevMode)
	if err != nil {
		return errors.Wrap(err, "fail init logger")
	}

	// init database config
	err = db.Init(config.DB, true)
	if err != nil {
		return errors.Wrap(err, "fail init db")
	}

	// check db connection
	err = db.Check(config.DB)
	if err != nil {
		return errors.Wrap(err, "fail check db connection")
	}

	// run migrations
	err = db.RunMigrations(config.DB)
	if err != nil {
		return errors.Wrap(err, "fail run db migrations")
	}

	// init file storage
	files.Init(config.CDNStorage.SelCDN, config.Params.AppDevMode)

	// init email sender
	messages.InitEmailSender(config.Notifications.Email)

	// init cache storage
	err = cache.Init(config.Cache)
	if err != nil {
		panic(err)
	}

	return nil
}
