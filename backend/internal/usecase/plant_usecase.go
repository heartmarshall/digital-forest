package usecase

import (
	"context"
	"time"

	// Убедись, что путь импорта соответствует имени твоего Go-модуля
	domain "github.com/heartmarshall/digital-forest/backend/internal/domain/plant"
)

// PlantUseCase - это конкретная реализация бизнес-логики для работы с "растениями".
// Он инкапсулирует все сценарии использования, связанные с этой сущностью.
type PlantUseCase struct {
	repo PlantRepository
}

// NewPlantUseCase - это конструктор для PlantUseCase.
// Он принимает PlantRepository в качестве зависимости и возвращает указатель
// на конкретную структуру *PlantUseCase.
func NewPlantUseCase(r PlantRepository) *PlantUseCase {
	return &PlantUseCase{repo: r}
}

// Create - это сценарий использования для создания нового растения.
// Он принимает базовые данные, формирует из них доменную модель Plant
// и передает ее в репозиторий для сохранения.
func (uc *PlantUseCase) Create(ctx context.Context, author, imageData string) (domain.Plant, error) {
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

func (uc *PlantUseCase) GetRandom(ctx context.Context, count int) ([]domain.Plant, error) {

	plants, err := uc.repo.GetRandom(ctx, count)
	if err != nil {
		return nil, err
	}

	return plants, nil
}
