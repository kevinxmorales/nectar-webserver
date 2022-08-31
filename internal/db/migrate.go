package db

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func (d *Database) MigrateDB() error {
	driver, err := postgres.WithInstance(d.Client.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("postgres.WithInstance in db.MigrateDB failed for %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver)
	if err != nil {
		return fmt.Errorf("migrate.NewWithDatabaseInstance in db.MigrateDB failed for %v", err)
	}
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("*migrate.Migrate.Up in db.MigrateDB failed for %v", err)
		}
	}
	log.Info("successfully migrated the database")
	return nil
}
