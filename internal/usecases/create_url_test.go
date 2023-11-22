package usecases

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/nikitades/url-shortener/internal/adapters"
	"github.com/nikitades/url-shortener/internal/usecases/mocks"
	"github.com/stretchr/testify/assert"
	tmock "github.com/stretchr/testify/mock"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestCreateUrl(t *testing.T) {
	db, dbmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create db mock: %s", err)
	}
	defer db.Close()

	exampleUrl := "https://www.sporks.org/blog/"
	exampleShortUrl := "abcdefg115"
	now := time.Now().In(time.UTC)

	dbmock.ExpectQuery("SELECT \\* FROM urls").WithArgs(exampleUrl).WillReturnError(sql.ErrNoRows)
	dbmock.ExpectPrepare("INSERT INTO urls")
	dbmock.ExpectQuery("INSERT INTO urls").WithArgs(exampleUrl, exampleShortUrl, now).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(418))

	urlRepo := adapters.NewPgsqlUrlRepo(sqlx.NewDb(db, "sqlmock"))

	urlgen := mocks.NewUrlGenerator(t)
	urlgen.On("Generate", tmock.AnythingOfType("int")).Return(exampleShortUrl)

	timeprov := mocks.NewTimeProvider(t)
	timeprov.On("Now").Return(now)

	createUrl := NewCreateUrlUsecase(
		urlRepo,
		urlgen,
		timeprov,
	)

	ctx := context.Background()

	result, err := createUrl(ctx, exampleUrl)
	assert.NoError(t, err)

	assert.Equal(t, 418, result.Id)
}

func TestGeneralErrorHandling(t *testing.T) {
	db, dbmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create db mock: %s", err)
	}
	defer db.Close()

	exampleUrl := "https://www.sporks.org/blog/"

	dbmock.ExpectQuery("SELECT \\* FROM urls").WillReturnError(sql.ErrConnDone)
	urlRepo := adapters.NewPgsqlUrlRepo(sqlx.NewDb(db, "sqlmock"))
	urlgen := mocks.NewUrlGenerator(t)
	timeprov := mocks.NewTimeProvider(t)

	createUrl := NewCreateUrlUsecase(
		urlRepo,
		urlgen,
		timeprov,
	)

	ctx := context.Background()

	_, err = createUrl(ctx, exampleUrl)
	assert.ErrorIs(t, err, sql.ErrConnDone)
}

func TestWhenUrlAlreadyExistsNoNewUrlIsCreated(t *testing.T) {
	db, dbmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create db mock: %s", err)
	}
	defer db.Close()

	exampleUrl := "https://www.sporks.org/blog/"

	dbmock.ExpectQuery("SELECT \\* FROM urls").WithArgs(exampleUrl).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(515))

	urlRepo := adapters.NewPgsqlUrlRepo(sqlx.NewDb(db, "sqlmock"))

	urlgen := mocks.NewUrlGenerator(t)

	timeprov := mocks.NewTimeProvider(t)

	createUrl := NewCreateUrlUsecase(
		urlRepo,
		urlgen,
		timeprov,
	)

	ctx := context.Background()

	_, err = createUrl(ctx, exampleUrl)
	assert.ErrorIs(t, err, BadRequestError)
	assert.ErrorContains(t, err, "url already exists")
}

func TestWhenHashCollisionHappensUsecaseRetries(t *testing.T) {
	db, dbmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create db mock: %s", err)
	}
	defer db.Close()

	exampleUrl := "https://www.sporks.org/blog/"
	exampleShortUrl := "abcdefg115"
	now := time.Now().In(time.UTC)

	dbmock.ExpectQuery("SELECT \\* FROM urls").WithArgs(exampleUrl).WillReturnError(sql.ErrNoRows)

	//first attempt
	dbmock.ExpectPrepare("INSERT INTO urls")
	uniqueViolationErr := pq.Error{
		Code: pq.ErrorCode("23505"),
	}
	dbmock.ExpectQuery("INSERT INTO urls").WithArgs(exampleUrl, exampleShortUrl, now).WillReturnError(&uniqueViolationErr)

	//second attempt
	dbmock.ExpectPrepare("INSERT INTO urls")
	dbmock.ExpectQuery("INSERT INTO urls").WithArgs(exampleUrl, exampleShortUrl, now).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(307))

	urlRepo := adapters.NewPgsqlUrlRepo(sqlx.NewDb(db, "sqlmock"))
	urlgen := mocks.NewUrlGenerator(t)
	urlgen.On("Generate", tmock.AnythingOfType("int")).Return(exampleShortUrl)

	timeprov := mocks.NewTimeProvider(t)
	timeprov.On("Now").Return(now)

	createUrl := NewCreateUrlUsecase(
		urlRepo,
		urlgen,
		timeprov,
	)

	ctx := context.Background()

	result, err := createUrl(ctx, exampleUrl)
	assert.NoError(t, err)
	assert.Equal(t, 307, result.Id)
}
