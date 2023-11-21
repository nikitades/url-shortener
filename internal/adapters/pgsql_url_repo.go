package adapters

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/nikitades/url-shortener/internal/domain"
)

type PgsqlUrlRepo struct {
	db *sqlx.DB
}

func (r *PgsqlUrlRepo) FindByOriginal(ctx context.Context, original string) (domain.Url, error) {
	row := r.db.QueryRowxContext(ctx, "SELECT * FROM urls WHERE source_url = $1", original)

	url := domain.Url{}
	err := row.StructScan(&url)

	if err == sql.ErrNoRows {
		return url, domain.NotFoundError
	}

	if err != nil {
		return url, err
	}

	return url, nil
}

func (r *PgsqlUrlRepo) FindByShort(ctx context.Context, short string) (domain.Url, error) {
	row := r.db.QueryRowxContext(ctx, "SELECT * FROM urls WHERE short_url = $1", short)

	url := domain.Url{}
	err := row.StructScan(&url)

	if err == sql.ErrNoRows {
		return url, domain.NotFoundError
	}

	if err != nil {
		return url, err
	}

	return url, nil
}

func (r *PgsqlUrlRepo) Create(ctx context.Context, original string, short string, createdAt time.Time) (domain.Url, error) {
	stmt, err := r.db.PrepareNamedContext(
		ctx,
		"INSERT INTO urls (source_url, short_url, created_at) VALUES (:source_url, :short_url, :created_at) RETURNING id",
	)
	url := domain.Url{}
	var id int

	if err != nil {
		return url, err
	}

	err = stmt.GetContext(ctx, &id, map[string]interface{}{
		"source_url": original,
		"short_url":  short,
		"created_at": createdAt,
	})

	if pqerr, ok := err.(*pq.Error); ok {
		if pqerr.Code.Name() == "unique_violation" {
			return domain.Url{}, domain.AlreadyExistsError
		}
		return domain.Url{}, pqerr
	}

	if err != nil {
		return url, err
	}

	url.Id = id
	url.SourceUrl = original
	url.ShortUrl = short
	url.CreatedAt = createdAt

	return url, nil
}

//TODO
/*
	- тесты

*/

func NewPgsqlUrlRepo(dbconn *sqlx.DB) *PgsqlUrlRepo {
	return &PgsqlUrlRepo{dbconn}
}

func NewPgsqlVisitRepo(dbconn *sqlx.DB) *PgsqlVisitRepo {
	return &PgsqlVisitRepo{dbconn}
}
