package usecases

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/nikitades/url-shortener/internal/adapters"
	"github.com/nikitades/url-shortener/internal/usecases/mocks"
	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetRedirect(t *testing.T) {
	db, dbmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create db mock: %s", err)
	}
	defer db.Close()

	urlCode := "abcdefgh123"
	userAgent := "testify-test"
	sourceUrl := "https://blank.org"

	urlRepo := adapters.NewPgsqlUrlRepo(sqlx.NewDb(db, "sqlmock"))
	visitsRepo := adapters.NewPgsqlVisitRepo(sqlx.NewDb(db, "sqlmock"))

	now := time.Now().In(time.UTC)
	timeprov := mocks.NewTimeProvider(t)
	timeprov.On("Now").Return(now)

	dbmock.ExpectQuery("SELECT \\* FROM urls").WithArgs(urlCode).WillReturnRows(sqlmock.NewRows([]string{"id", "source_url", "short_url"}).AddRow(506, sourceUrl, urlCode))
	dbmock.ExpectPrepare("INSERT INTO visits")
	dbmock.ExpectQuery("INSERT INTO visits").WithArgs(userAgent, 506, sourceUrl, urlCode, now).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(704))

	ctx := context.Background()

	getRedirect := NewGetRedirectUsecase(urlRepo, visitsRepo, timeprov)

	redirect, err := getRedirect(ctx, urlCode, userAgent)
	assert.NoError(t, err)
	assert.Equal(t, sourceUrl, redirect)
}

func TestWhenRedirectIsNotFoundProperErrorIsReturned(t *testing.T) {
	db, dbmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create db mock: %s", err)
	}
	defer db.Close()

	urlCode := "abcdefgh123"

	urlRepo := adapters.NewPgsqlUrlRepo(sqlx.NewDb(db, "sqlmock"))
	visitsRepo := adapters.NewPgsqlVisitRepo(sqlx.NewDb(db, "sqlmock"))
	timeprov := mocks.NewTimeProvider(t)

	dbmock.ExpectQuery("SELECT \\* FROM urls").WithArgs(urlCode).WillReturnError(sql.ErrNoRows)

	ctx := context.Background()

	getRedirect := NewGetRedirectUsecase(urlRepo, visitsRepo, timeprov)

	_, err = getRedirect(ctx, urlCode, "privet")
	assert.ErrorIs(t, err, NotFoundError)
}

func TestWhenFailedToFetchUrlErrorIsGiven(t *testing.T) {
	db, dbmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create db mock: %s", err)
	}
	defer db.Close()

	urlCode := "abcdefgh123"

	urlRepo := adapters.NewPgsqlUrlRepo(sqlx.NewDb(db, "sqlmock"))
	visitsRepo := adapters.NewPgsqlVisitRepo(sqlx.NewDb(db, "sqlmock"))
	timeprov := mocks.NewTimeProvider(t)

	dbmock.ExpectQuery("SELECT \\* FROM urls").WithArgs(urlCode).WillReturnError(sql.ErrConnDone)

	ctx := context.Background()

	getRedirect := NewGetRedirectUsecase(urlRepo, visitsRepo, timeprov)

	_, err = getRedirect(ctx, urlCode, "privet")
	assert.ErrorIs(t, err, sql.ErrConnDone)
}

func TestWhenFailedToCreateVisitErrorIsGiven(t *testing.T) {
	db, dbmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create db mock: %s", err)
	}
	defer db.Close()

	urlCode := "abcdefgh123"

	urlRepo := adapters.NewPgsqlUrlRepo(sqlx.NewDb(db, "sqlmock"))
	visitsRepo := adapters.NewPgsqlVisitRepo(sqlx.NewDb(db, "sqlmock"))

	timeprov := mocks.NewTimeProvider(t)
	timeprov.On("Now").Return(time.Now().In(time.UTC))

	dbmock.ExpectQuery("SELECT \\* FROM urls").WithArgs(urlCode).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(61))
	dbmock.ExpectPrepare("INSERT INTO visits")
	dbmock.ExpectQuery("INSERT INTO visits").WillReturnError(sql.ErrConnDone)

	ctx := context.Background()

	getRedirect := NewGetRedirectUsecase(urlRepo, visitsRepo, timeprov)

	_, err = getRedirect(ctx, urlCode, "privet")
	assert.ErrorIs(t, err, sql.ErrConnDone)
}
