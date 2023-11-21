package adapters

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/nikitades/url-shortener/internal/domain"
)

type PgsqlVisitRepo struct {
	db *sqlx.DB
}

func (r *PgsqlVisitRepo) Create(ctx context.Context, urlId int, urlSource string, urlCode string, userAgent string, createdAt time.Time) (domain.Visit, error) {
	stmt, err := r.db.PrepareNamedContext(
		ctx,
		"INSERT INTO visits (user_agent, url_id, url_source, url_code, created_at) VALUES (:user_agent, :url_id, :url_source, :url_code, :created_at) RETURNING id",
	)
	visit := domain.Visit{}
	var id int

	if err != nil {
		return visit, err
	}

	err = stmt.GetContext(ctx, &id, map[string]interface{}{
		"user_agent": userAgent,
		"url_id":     urlId,
		"url_source": urlSource,
		"url_code":   urlCode,
		"created_at": createdAt,
	})

	if err != nil {
		return visit, err
	}

	visit.Id = id
	visit.UserAgent = userAgent
	visit.UrlId = urlId
	visit.UrlSource = urlSource
	visit.UrlCode = urlCode
	visit.CreatedAt = createdAt

	return visit, nil
}

func (r *PgsqlVisitRepo) FindByUrlCode(ctx context.Context, urlCode string) ([]domain.Visit, error) {
	visits := []domain.Visit{}
	err := r.db.SelectContext(ctx, &visits, "SELECT * FROM visits WHERE url_code = $1", urlCode)
	
	if err != nil {
		return visits, err
	}

	return visits, nil
}
