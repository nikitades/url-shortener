package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/nikitades/url-shortener/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWhenBadPayloadProvidedBadRequestIsGiven(t *testing.T) {
	setupTest(t)

	r := httptest.NewRequest("POST", "/api/v1/url", strings.NewReader("{\"bad json\"}"))
	rr := httptest.NewRecorder()
	_api.r.ServeHTTP(rr, r)

	assert.Equal(t, http.StatusBadRequest, rr.Result().StatusCode)
}

func TestWhenUrlAlreadyExistsBadRequestIsGiven(t *testing.T) {
	setupTest(t)

	sourceUrl := "https://www.sporks.org/blog/"

	payload, err := json.Marshal(map[string]string{
		"source": sourceUrl,
	})
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest("POST", "/api/v1/url", bytes.NewReader(payload))

	urlrepo.On("FindByOriginal", mock.Anything, sourceUrl).Return(domain.Url{}, nil)

	rr := httptest.NewRecorder()
	_api.r.ServeHTTP(rr, r)

	assert.Equal(t, http.StatusBadRequest, rr.Result().StatusCode)
	txt, err := io.ReadAll(rr.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(txt), "already exists")
}

func TestWhenUnknownErrorHappensInternalServerErrorIsGiven(t *testing.T) {
	setupTest(t)

	sourceUrl := "https://www.sporks.org/blog/"

	payload, err := json.Marshal(map[string]string{
		"source": sourceUrl,
	})
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest("POST", "/api/v1/url", bytes.NewReader(payload))

	urlrepo.On("FindByOriginal", mock.Anything, sourceUrl).Return(domain.Url{}, sql.ErrConnDone)

	rr := httptest.NewRecorder()
	_api.r.ServeHTTP(rr, r)

	assert.Equal(t, http.StatusInternalServerError, rr.Result().StatusCode)
	txt, err := io.ReadAll(rr.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(txt), "internal error")
}

func TestNewUrlIsCreated(t *testing.T) {
	setupTest(t)

	sourceUrl := "https://www.sporks.org/blog/"
	sampleCode := "abcdefg1234"
	now := time.Now().In(time.UTC)
	expectedUrl := domain.Url{
		Id:        714,
		SourceUrl: sourceUrl,
		ShortUrl:  sampleCode,
		CreatedAt: now,
	}

	payload, err := json.Marshal(map[string]string{
		"source": sourceUrl,
	})
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest("POST", "/api/v1/url", bytes.NewReader(payload))

	urlrepo.On("FindByOriginal", mock.Anything, sourceUrl).Return(domain.Url{}, domain.NotFoundError)
	urlgen.On("Generate", 12).Return(sampleCode)
	timeprov.On("Now").Return(now)
	urlrepo.On("Create", mock.Anything, sourceUrl, sampleCode, now).Return(expectedUrl, nil)

	rr := httptest.NewRecorder()
	_api.r.ServeHTTP(rr, r)

	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
	createdUrl := domain.Url{}
	json.NewDecoder(rr.Body).Decode(&createdUrl)
	assert.Equal(t, expectedUrl, createdUrl)
}
