package migration

import (
	"database/sql"
	"errors"
	"github.com/EestiChameleon/GOphermart/internal/app/cfg"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var (
	m *migrate.Migrate
)

func UpGophermartStorage() error {
	err := m.Up()
	if !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

func DownGophermartStorage() error {
	return m.Down()
}

func MigrateInitConnect() error {
	conn, err := sql.Open("postgres",
		cfg.Envs.DatabaseURI)

	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		return err
	}

	db, err := migrate.NewWithDatabaseInstance(
		"file://migration/sqlscripts/",
		"postgres", driver)
	if err != nil {
		return err
	}

	m = db
	return nil
}

func MigrateCloseConnect() (error, error) {
	return m.Close()
}
