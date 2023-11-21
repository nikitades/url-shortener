package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/nikitades/url-shortener/internal/adapters"
	"github.com/nikitades/url-shortener/internal/domain"
	"github.com/nikitades/url-shortener/internal/usecases"
)

type api struct {
	r        *chi.Mux
	usecases *usecases.Usecases
}

func (api *api) init() {
	api.r.Get("/api/v1/health", api.healthcheck)
	api.r.Post("/api/v1/url", api.createNewUrl)
	api.r.Get("/{code}", api.redirect)
	api.r.Get("/api/v1/url/{code}/stats", api.getStats)
}

var createDependencies = func(sqlconnstr string) (
	*sqlx.DB,
	domain.UrlRepository,
	domain.VisitRepository,
	usecases.UrlGenerator,
	*usecases.TimeProvider,
) {
	dbconn, err := adapters.NewPgsqlConn(sqlconnstr)
	if err != nil {
		log.Fatal(err)
	}

	urlRepo := adapters.NewPgsqlUrlRepo(dbconn)
	visitRepo := adapters.NewPgsqlVisitRepo(dbconn)
	urlgen := adapters.NewRandBytesUrlGenerator()
	timeprov := &usecases.TimeProvider{}

	return dbconn, urlRepo, visitRepo, urlgen, timeprov
}

func newApi(port string, sqlconnstr string) *api {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(contentTypeMiddleware)

	usecases, err := usecases.InitUsecases(
		createDependencies(sqlconnstr),
	)

	if err != nil {
		log.Fatal(err)
	}

	return &api{
		r,
		usecases,
	}
}

func Start(port string, sqlconnstr string) {
	api := newApi(port, sqlconnstr)
	api.init()
	log.Printf("started at :%s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), api.r)
}
