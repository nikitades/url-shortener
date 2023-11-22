package usecases

import (
	"context"

	"github.com/nikitades/url-shortener/internal/domain"
)

type GetVisitsByUrlCodeUsecase func(ctx context.Context, urlCode string) ([]domain.Visit, error)

func NewGetVisitsByUrlCodeUsecase(visitRepo domain.VisitRepository, urlRepo domain.UrlRepository) GetVisitsByUrlCodeUsecase {
	return func(ctx context.Context, urlCode string) ([]domain.Visit, error) {
		_, err := urlRepo.FindByShort(ctx, urlCode)

		if err == domain.NotFoundError {
			return []domain.Visit{}, NotFoundError
		}

		if err != nil {
			return []domain.Visit{}, err
		}

		visits, err := visitRepo.FindByUrlCode(ctx, urlCode)

		if err != nil {
			return []domain.Visit{}, err
		}

		return visits, nil
	}
}
