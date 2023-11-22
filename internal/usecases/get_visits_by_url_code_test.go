package usecases

import (
	"context"
	"database/sql"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/nikitades/url-shortener/internal/adapters"
	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetVisitsByUrlCode(t *testing.T) {
	db, dbmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create db mock: %s", err)
	}
	defer db.Close()

	urlcode := "abcdef123"

	visitRepo := adapters.NewPgsqlVisitRepo(sqlx.NewDb(db, "sqlmock"))
	urlRepo := adapters.NewPgsqlUrlRepo(sqlx.NewDb(db, "sqlmock"))
	dbmock.ExpectQuery("SELECT \\* FROM urls").WithArgs(urlcode).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(716))
	dbmock.ExpectQuery("SELECT \\* FROM visits").WithArgs(urlcode).WillReturnRows(sqlmock.NewRows([]string{"id"}).
		AddRow(808).
		AddRow(501).
		AddRow(109).
		AddRow(500),
	)
	ctx := context.Background()

	getVisitsByUrlCode := NewGetVisitsByUrlCodeUsecase(visitRepo, urlRepo)
	visits, err := getVisitsByUrlCode(ctx, urlcode)

	assert.NoError(t, err)
	assert.Len(t, visits, 4)
}

func TestWhenErrorHappensErrorIsGiven(t *testing.T) {
	db, dbmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create db mock: %s", err)
	}
	defer db.Close()

	urlcode := "abcdef123"

	visitRepo := adapters.NewPgsqlVisitRepo(sqlx.NewDb(db, "sqlmock"))
	urlRepo := adapters.NewPgsqlUrlRepo(sqlx.NewDb(db, "sqlmock"))
	dbmock.ExpectQuery("SELECT \\* FROM urls").WithArgs(urlcode).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(716))
	dbmock.ExpectQuery("SELECT \\* FROM visits").WithArgs(urlcode).WillReturnError(sql.ErrConnDone)
	ctx := context.Background()

	getVisitsByUrlCode := NewGetVisitsByUrlCodeUsecase(visitRepo, urlRepo)
	_, err = getVisitsByUrlCode(ctx, urlcode)

	assert.ErrorIs(t, err, sql.ErrConnDone)
}
