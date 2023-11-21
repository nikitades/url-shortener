package usecases

import (
	"context"

	"github.com/nikitades/url-shortener/internal/domain"
)

type GetVisitsByUrlCodeUsecase func(context.Context, string) ([]domain.Visit, error)

func NewGetVisitsByUrlCodeUsecase(visitRepo domain.VisitRepository) GetVisitsByUrlCodeUsecase {
	return func(ctx context.Context, urlCode string) ([]domain.Visit, error) {
		visits, err := visitRepo.FindByUrlCode(ctx, urlCode)

		if err != nil {
			return []domain.Visit{}, err
		}

		return visits, nil
	}
}
