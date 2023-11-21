package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/nikitades/url-shortener/internal/usecases"
)

type newUrlRequest struct {
	Source string
}

func (api *api) createNewUrl(w http.ResponseWriter, r *http.Request) {
	urlRequest := newUrlRequest{}
	err := json.NewDecoder(r.Body).Decode(&urlRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(HttpError{"bad request"})
		return
	}

	url, err := api.usecases.CreateUrl(r.Context(), urlRequest.Source)
	if errors.Is(err, usecases.BadRequestError) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(HttpError{fmt.Sprintf("%s", err)})
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(HttpError{"internal error"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(url)
}
