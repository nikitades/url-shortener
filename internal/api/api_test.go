package api

import (
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	domainmocks "github.com/nikitades/url-shortener/internal/domain/mocks"
	"github.com/nikitades/url-shortener/internal/usecases"
	usecasemocks "github.com/nikitades/url-shortener/internal/usecases/mocks"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func setupTest(t *testing.T) {
	r := chi.NewRouter()
	r.Use(contentTypeMiddleware)

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create db mock: %s", err)
	}
	defer db.Close()

	dbconn := sqlx.NewDb(db, "sqlmock")
	urlrepo = domainmocks.NewUrlRepository(t)
	visitrepo = domainmocks.NewVisitRepository(t)
	urlgen = usecasemocks.NewUrlGenerator(t)
	timeprov = usecasemocks.NewTimeProvider(t)

	usecases, err := usecases.InitUsecases(
		dbconn,
		urlrepo,
		visitrepo,
		urlgen,
		timeprov,
	)

	if err != nil {
		t.Fatal(err)
	}

	_api = &api{
		r,
		usecases,
	}

	_api.init()
}
