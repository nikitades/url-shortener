package usecases

import "time"

//go:generate mockery --name TimeProvider
type TimeProvider interface {
	Now() time.Time
}