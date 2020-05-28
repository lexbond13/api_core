package migrations

import (
	"database/sql"
	"fmt"
	"github.com/lexbond13/api_core/module/logger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"time"
)

type Postgres struct {
	db  *sql.DB
	log logger.ILogger
}

func NewPostgres(log logger.ILogger, u string, pathToMigrations string) (*Postgres, error) {

	db, err := sql.Open("postgres", u)
	if err != nil {
		return nil, fmt.Errorf("unable to init id: %s", err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	var driver database.Driver

	for i := 0; i < 5; i++ {
		driver, err = postgres.WithInstance(db, &postgres.Config{})
		if err == nil {
			break
		}
		log.Info(err)
		time.Sleep(time.Second * 3)
	}

	if err != nil {
		return nil, fmt.Errorf("migrations driver creation error: %s", err)
	}

	m, err := migrate.NewWithDatabaseInstance(pathToMigrations, "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("migrations instance creation error: %s", err)
	}


	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("migrations error: %s", err)
	}

	psql := &Postgres{
		db:  db,
		log: log,
	}

	return psql, nil
}

func (pg *Postgres) Ping() error {
	// pg.db.Ping doesn't return error even when database is not alive
	_, err := pg.db.Exec(`SELECT 1`)
	return err
}

func (pg *Postgres) Close() error {
	return pg.db.Close()
}
