package usecases

import (
	"time"

	"github.com/heptiolabs/healthcheck"
	"github.com/jmoiron/sqlx"
)

type HealthCheckUsecase func() error

func NewHealthCheckUsecase(db *sqlx.DB) HealthCheckUsecase {
	return func() error {
		dbcheck := healthcheck.DatabasePingCheck(db.DB, time.Second)
		if err := dbcheck(); err != nil {
			return err
		}
		return nil
	}
}
