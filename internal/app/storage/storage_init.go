package storage

import (
	"context"
	"errors"
	"github.com/EestiChameleon/GOphermart/internal/app/cfg"
	"github.com/EestiChameleon/GOphermart/migration"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	Pool        DBStorage
	ErrNotFound = errors.New("no records found")
)

type DBStorage struct {
	ID int // user_id of the current session - obtained from cookie via auth MW
	DB *pgxpool.Pool
}

func InitConnection() error {
	//create tables if it doesn't exist
	if err := migration.MigrateInitConnect(); err != nil {
		return err
	}

	if err := migration.UpGophermartStorage(); err != nil {
		return err
	}
	migration.MigrateCloseConnect()

	return connectToDB()
}

func Shutdown() {
	Pool.DB.Close()

}

func connectToDB() error {
	conn, err := pgxpool.Connect(context.Background(), cfg.Envs.DatabaseURI)
	if err != nil {
		return err
	}

	Pool = DBStorage{DB: conn}
	return nil
}
