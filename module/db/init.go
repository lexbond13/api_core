package db

import (
	"fmt"
	"github.com/lexbond13/api_core/config"
	"github.com/lexbond13/api_core/module/db/migrations"
	"github.com/lexbond13/api_core/module/logger"
	"github.com/go-pg/pg/v9"
)

var (
	connection *pg.DB
)

// Init
func Init(config *config.DB, debug bool) error {
	connection = pg.Connect(&pg.Options{
		User:     config.Username,
		Password: config.Password,
		Database: config.Database,
		PoolSize: config.PoolSize,
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
	})

	if debug {
		connection.AddQueryHook(&queryDebug{})
	}

	return nil
}

// Check connect try with current setting
func Check(config *config.DB) error {
	db := pg.Connect(&pg.Options{
		User:     config.Username,
		Password: config.Password,
		Database: config.Database,
		PoolSize: config.PoolSize,
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
	})
	defer func() {
		if err := db.Close(); err != nil {
			logger.Log.Error(err)
		}
	}()

	var n int
	_, err := db.QueryOne(pg.Scan(&n), "SELECT 1")
	if err != nil {
		return fmt.Errorf("can't connect to PSQL: %v", err)
	}
	return nil
}

func RunMigrations(config *config.DB) error {
	dbConnectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", config.Username, config.Password, config.Host, config.Port, config.Database)
	psql, err := migrations.NewPostgres(logger.Log, dbConnectionString, "file://module/db/migrations/")
	if err != nil {
		return err
	}
	defer func() {
		if err := psql.Close(); err != nil {
			logger.Log.Error(err)
		}
	}()

	return nil
}

// GetConnection
func GetConnection() *pg.DB {
	return connection
}
