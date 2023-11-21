package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nikitades/url-shortener/internal/usecases"
)

type StatsResponse struct {
	Stats map[string]int `json:"stats"`
}

func (api *api) getStats(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	visits, err := api.usecases.GetVisitsByUrlCode(r.Context(), code)

	if err == usecases.NotFoundError {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(HttpError{fmt.Sprintf("%s", err)})
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(HttpError{fmt.Sprintf("%s", err)})
		return
	}

	stats := make(map[string]int)
	for _, v := range visits {
		if _, ok := stats[v.UrlSource]; !ok {
			stats[v.UrlSource] = 0
		}
		stats[v.UrlSource]++
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(StatsResponse{stats})
}
