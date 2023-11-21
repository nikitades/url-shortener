package usecases

import (
	"context"

	"github.com/nikitades/url-shortener/internal/domain"
)

type GetRedirectUsecase func(context.Context, string, string) (string, error)

func NewGetRedirectUsecase(
	urlRepo domain.UrlRepository,
	visitRepo domain.VisitRepository,
	timeprov *TimeProvider,
) GetRedirectUsecase {

	return func(ctx context.Context, code string, userAgent string) (string, error) {
		url, err := urlRepo.FindByShort(ctx, code)

		if err == domain.NotFoundError {
			return "", NotFoundError
		}

		if err != nil {
			return "", err
		}

		_, err = visitRepo.Create(ctx, url.Id, url.SourceUrl, url.ShortUrl, userAgent, timeprov.Now())

		if err != nil {
			return "", err
		}

		return url.SourceUrl, nil
	}
}
