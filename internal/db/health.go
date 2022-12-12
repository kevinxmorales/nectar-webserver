package db

import (
	"context"
	log "github.com/sirupsen/logrus"
)

func (d *Database) CheckDbHealth(ctx context.Context) error {
	log.Info("Pinging database...")
	return d.Client.PingContext(ctx)
}
