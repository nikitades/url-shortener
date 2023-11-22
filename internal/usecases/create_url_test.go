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

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	dbmock.ExpectQuery("SELECT \\* FROM urls").WithArgs(exampleUrl).WillReturnError(sql.ErrNoRows)
	dbmock.ExpectPrepare("INSERT INTO urls")
	dbmock.ExpectQuery("INSERT INTO urls").WithArgs(exampleUrl, exampleShortUrl, now).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(418))

	urlRepo := adapters.NewPgsqlUrlRepo(sqlxDB)

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
