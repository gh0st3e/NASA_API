package db

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/gh0st3e/NASA_API/internal/config"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
)

func InitPSQL(cfg config.PSQLDatabase) (*sql.DB, error) {
	db, err := sql.Open(cfg.Driver, cfg.Address)
	if err != nil {
		return nil, fmt.Errorf("couldn't open PSQL: %s", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("couldn't ping PSQl: %s", err)
	}

	logrus.Info("Ping PSQL - OK!")

	err = migrateTables(db)
	if err != nil {
		return nil, fmt.Errorf("couldn't migrate tables: %s", err)
	}

	return db, nil
}

//go:embed migrations/*.sql
var embedTablesMigrations embed.FS

func migrateTables(db *sql.DB) error {

	goose.SetBaseFS(embedTablesMigrations)

	err := goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	err = goose.Up(db, "migrations", goose.WithAllowMissing())

	return err
}
