package domain

import (
	"context"
	"time"
)

//go:generate mockery --name VisitRepository
type VisitRepository interface {
	FindByUrlCode(ctx context.Context, urlCode string) ([]Visit, error)
	Create(
		ctx context.Context,
		urlId int,
		urlSource string,
		urlCode string,
		userAgent string,
		createdAt time.Time,
	) (Visit, error)
}
