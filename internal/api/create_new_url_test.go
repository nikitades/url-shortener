package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nikitades/url-shortener/internal/domain"
	domainmocks "github.com/nikitades/url-shortener/internal/domain/mocks"
	usecasemocks "github.com/nikitades/url-shortener/internal/usecases/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var urlrepo *domainmocks.UrlRepository
var visitrepo *domainmocks.VisitRepository
var urlgen *usecasemocks.UrlGenerator
var timeprov *usecasemocks.TimeProvider

var _api *api

func TestWhenBadPayloadProvidedBadRequestIsGiven(t *testing.T) {
	setupTest(t)

	r, err := http.NewRequest("POST", "/api/v1/url", strings.NewReader("{\"bad json\"}"))
	if err != nil {
		t.Fatal(err)
	}

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

	r, err := http.NewRequest("POST", "/api/v1/url", bytes.NewReader(payload))

	urlrepo.On("FindByOriginal", mock.Anything, sourceUrl).Return(domain.Url{}, nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	_api.r.ServeHTTP(rr, r)

	assert.Equal(t, http.StatusBadRequest, rr.Result().StatusCode)
	txt, err := io.ReadAll(rr.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(txt), "already exists")
}
