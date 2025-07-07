package database

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var (
	SX *sqlx.DB
)

type SSLMode string

const (
	SSLModeDisable SSLMode = "disable"
	SSLModeEnable  SSLMode = "enable"
)

func ConnectAndMigrate(host, port, databaseName, user, password string, ssl SSLMode) error {
	logrus.Info("Connecting to database...")
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s  sslmode=%s", host, port, user, password, databaseName, ssl)
	DB, dbError := sqlx.Connect("postgres", connectionString)
	if dbError != nil {
		logrus.Errorf("failed in connecting...")
		return dbError
	}
	dbError = DB.Ping()
	if dbError != nil {
		logrus.Errorf("failed in pinging...")
		return dbError
	}
	logrus.Info("Successfully Connected Database")
	SX = DB
	//call migration func
	return migrateUpAndDown(DB)
}

func ShutdownDB() error {
	logrus.Info("Shutting down database")
	return SX.Close()
}

func migrateUpAndDown(db *sqlx.DB) error {
	logrus.Info("Migrating database...")
	dbDriver, dbError := postgres.WithInstance(db.DB, &postgres.Config{})
	if dbError != nil {
		logrus.Errorf("failed in making db instance...")
		return dbError
	}
	path := "file://database/migrations/"
	dbName := "postgres"
	mig, migError := migrate.NewWithDatabaseInstance(path, dbName, dbDriver)
	if migError != nil {
		logrus.Errorf("failed in migrating...")
		return migError
	}
	if migError = mig.Up(); migError != nil && !errors.Is(migError, migrate.ErrNoChange) {
		logrus.Errorf("failed in .up file migration...")
		return migError
	}
	if migError = mig.Down(); migError != nil && !errors.Is(migError, migrate.ErrNoChange) {
		logrus.Info("failed in .down file migration...")
		return migError
	}
	logrus.Info("Successfully migrated database")
	return nil
}
