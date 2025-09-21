package create

import (
	"context"
	"time"

	domain "github.com/heartmarshall/digital-forest/backend/internal/domain/plant"
)

// PlantRepository определяет контракт для слоя данных.
type PlantRepository interface {
	Create(ctx context.Context, plant domain.Plant) (domain.Plant, error)
}

// CreateUseCase - это конкретная реализация бизнес-логики для создания растения.
type CreateUseCase struct {
	repo PlantRepository
}

// NewCreateUseCase - конструктор для CreateUseCase.
func NewCreateUseCase(r PlantRepository) *CreateUseCase {
	return &CreateUseCase{repo: r}
}

// Create - сценарий использования для создания нового растения.
func (uc *CreateUseCase) Create(ctx context.Context, author, imageData string) (domain.Plant, error) {
	// Здесь в будущем могла бы быть бизнес-валидация.
	// Например, проверка imageData на корректность формата,
	// или проверка имени автора на наличие в черном списке.

	plant := domain.Plant{
		Author:    author,
		ImageData: imageData,
		CreatedAt: time.Now().UTC(),
	}

	createdPlant, err := uc.repo.Create(ctx, plant)
	if err != nil {
		return domain.Plant{}, err
	}
	return createdPlant, nil
}
