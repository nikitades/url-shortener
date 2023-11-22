package usecases

import (
	"context"
	"fmt"

	"github.com/nikitades/url-shortener/internal/domain"
)

type CreateUrlUsecase func(context.Context, string) (domain.Url, error)

func NewCreateUrlUsecase(repo domain.UrlRepository, urlgen UrlGenerator, timeprov TimeProvider) CreateUrlUsecase {

	return func(ctx context.Context, s string) (domain.Url, error) {
		_, err := repo.FindByOriginal(ctx, s)
		if err != nil && err != domain.NotFoundError {
			return domain.Url{}, err
		}

		if err == nil {
			return domain.Url{}, fmt.Errorf("%w: url already exists", BadRequestError)
		}

		for {
			url, err := repo.Create(ctx, s, urlgen.Generate(12), timeprov.Now())

			if err == domain.AlreadyExistsError {
				continue //try again
			}

			if err != nil {
				return domain.Url{}, err
			}
			return url, nil
		}
	}
}
