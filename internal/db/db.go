package db

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"os"
)

type Database struct {
	Client *sqlx.DB
}

func NewDatabase() (*Database, error) {
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_DB"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("SSL_MODE"))
	dbConn, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("sqlx.Connect in NewDatabase failed for %w", err)
	}
	database := Database{Client: dbConn}
	return &database, nil
}

func (d *Database) Ping(ctx context.Context) error {
	return d.Client.PingContext(ctx)
}

type SqlRows interface {
	Close() error
	Next() bool
	Scan(dest ...any) error
}

func closeDbRows(rows SqlRows, query string) {
	if err := rows.Close(); err != nil {
		log.Errorf("FAILED to close rows from query %s", query)
	}
}

func convertList[T any, U any](inputList []T, convertFunc func(T) U) []U {
	outputList := make([]U, len(inputList))
	for i, v := range inputList {
		outputList[i] = convertFunc(v)
	}
	return outputList
}
