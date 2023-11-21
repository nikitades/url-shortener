package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nikitades/url-shortener/internal/usecases"
)

func (api *api) redirect(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	redirect, err := api.usecases.Redirect(r.Context(), code, r.UserAgent())
	if errors.Is(err, usecases.NotFoundError) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(HttpError{fmt.Sprintf("%s", err)})
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(HttpError{fmt.Sprintf("%s", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"redirect": redirect,
	})
}
