package domain

import (
	"context"
	"time"
)

//go:generate mockery --name UrlRepository
type UrlRepository interface {
	FindByOriginal(ctx context.Context, original string) (Url, error)
	FindByShort(ctx context.Context, short string) (Url, error)
	Create(ctx context.Context, original string, short string, createdAt time.Time) (Url, error)
}
