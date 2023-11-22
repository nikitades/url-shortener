package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nikitades/url-shortener/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWhenNoUrlNotFoundErrorIsGiven(t *testing.T) {
	setupTest(t)

	urlCode := "abcdefg1234"

	urlrepo.On("FindByShort", mock.Anything, urlCode).Return(domain.Url{}, domain.NotFoundError)

	r := httptest.NewRequest("GET", "/"+urlCode, nil)
	rr := httptest.NewRecorder()

	_api.r.ServeHTTP(rr, r)

	assert.Equal(t, http.StatusNotFound, rr.Result().StatusCode)
}

func TestWhenUrlIsFoundRedirectHappens(t *testing.T) {
	setupTest(t)

	urlCode := "abcdefg1234"
	urlSource := "https://sporks.org/blog/"
	now := time.Now().In(time.UTC)
	expectedUrl := domain.Url{Id: 509, SourceUrl: urlSource, ShortUrl: urlCode, CreatedAt: now}

	urlrepo.On("FindByShort", mock.Anything, urlCode).Return(expectedUrl, nil)
	visitrepo.On("Create", mock.Anything, expectedUrl.Id, expectedUrl.SourceUrl, expectedUrl.ShortUrl, "euphoria", now).Return(domain.Visit{}, nil)
	timeprov.On("Now").Return(now)

	r := httptest.NewRequest("GET", "/"+urlCode, nil)
	r.Header.Add("user-agent", "euphoria")
	rr := httptest.NewRecorder()

	_api.r.ServeHTTP(rr, r)

	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
	expected, err := json.Marshal(map[string]string{"redirect": expectedUrl.SourceUrl})
	assert.NoError(t, err)
	actual, err := io.ReadAll(rr.Body)
	assert.NoError(t, err)
	assert.Equal(t, string(expected)+"\n", string(actual))
}
