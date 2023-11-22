package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nikitades/url-shortener/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWhenNoUrlNotFoundIsGiven(t *testing.T) {
	setupTest(t)

	urlCode := "abcdefg1234"

	urlrepo.On("FindByShort", mock.Anything, urlCode).Return(domain.Url{}, domain.NotFoundError)

	r := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/url/%s/stats", urlCode), nil)
	rr := httptest.NewRecorder()

	_api.r.ServeHTTP(rr, r)
	assert.Equal(t, http.StatusNotFound, rr.Result().StatusCode)
}

func TestWhenUnknownErrorHappensInternalServerErrorIsGivenAtGetStats(t *testing.T) {
	setupTest(t)

	urlCode := "abcdefg1234"

	urlrepo.On("FindByShort", mock.Anything, urlCode).Return(domain.Url{}, sql.ErrConnDone)

	r := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/url/%s/stats", urlCode), nil)
	rr := httptest.NewRecorder()

	_api.r.ServeHTTP(rr, r)
	assert.Equal(t, http.StatusInternalServerError, rr.Result().StatusCode)
}

func TestVisitsAreFound(t *testing.T) {
	setupTest(t)

	urlSource := "https://sporks.org/blog/"
	urlCode := "abcdefg1234"

	urlrepo.On("FindByShort", mock.Anything, urlCode).Return(domain.Url{}, nil)
	visitrepo.On("FindByUrlCode", mock.Anything, urlCode).Return(
		[]domain.Visit{{UrlSource: urlSource}, {UrlSource: urlSource}, {UrlSource: urlSource}, {UrlSource: urlSource}, {UrlSource: urlSource}},
		nil,
	)

	r := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/url/%s/stats", urlCode), nil)
	rr := httptest.NewRecorder()

	_api.r.ServeHTTP(rr, r)
	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
	expected, err := json.Marshal(map[string]map[string]int{"stats": {
		urlSource: 5,
	}})
	assert.NoError(t, err)
	actual, err := io.ReadAll(rr.Body)
	assert.NoError(t, err)
	assert.Equal(t, string(expected)+"\n", string(actual))
}
