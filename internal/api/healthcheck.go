package api

import (
	"fmt"
	"net/http"
)

func (api *api) healthcheck(w http.ResponseWriter, r *http.Request) {
	err := api.usecases.HealthCheck()
	if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(fmt.Sprintf("failed to init app: %s", err)))
}
