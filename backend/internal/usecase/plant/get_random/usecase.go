package get_random

import (
	"context"

	domain "github.com/heartmarshall/digital-forest/backend/internal/domain/plant"
)

// PlantRepository определяет контракт для слоя данных.
type PlantRepository interface {
	GetRandom(ctx context.Context, count int) ([]domain.Plant, error)
}

// GetRandomUseCase - это конкретная реализация бизнес-логики для получения случайных растений.
type GetRandomUseCase struct {
	repo PlantRepository
}

// NewGetRandomUseCase - конструктор для GetRandomUseCase.
func NewGetRandomUseCase(r PlantRepository) *GetRandomUseCase {
	return &GetRandomUseCase{repo: r}
}

// GetRandom - сценарий использования для получения случайных растений.
func (uc *GetRandomUseCase) GetRandom(ctx context.Context, count int) ([]domain.Plant, error) {
	plants, err := uc.repo.GetRandom(ctx, count)
	if err != nil {
		return nil, err
	}

	return plants, nil
}
