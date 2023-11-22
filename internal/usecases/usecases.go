package usecases

import (
	"github.com/jmoiron/sqlx"
	"github.com/nikitades/url-shortener/internal/domain"
)

type Usecases struct {
	HealthCheck        HealthCheckUsecase
	CreateUrl          CreateUrlUsecase
	Redirect           GetRedirectUsecase
	GetVisitsByUrlCode GetVisitsByUrlCodeUsecase
}

func InitUsecases(
	dbconn *sqlx.DB,
	urlRepo domain.UrlRepository,
	visitRepo domain.VisitRepository,
	urlgen UrlGenerator,
	timeprov TimeProvider,
) (*Usecases, error) {
	return &Usecases{
		HealthCheck:        NewHealthCheckUsecase(dbconn),
		CreateUrl:          NewCreateUrlUsecase(urlRepo, urlgen, timeprov),
		Redirect:           NewGetRedirectUsecase(urlRepo, visitRepo, timeprov),
		GetVisitsByUrlCode: NewGetVisitsByUrlCodeUsecase(visitRepo, urlRepo),
	}, nil
}
