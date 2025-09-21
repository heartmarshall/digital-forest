package usecase

import (
	"context"

	// Убедись, что путь импорта соответствует имени твоего Go-модуля
	domain "github.com/heartmarshall/digital-forest/backend/internal/domain/plant"
)

// PlantRepository определяет контракт ДЛЯ СЛОЯ ДАННЫХ, который ДОЛЖЕН БЫТЬ реализован.
// Слой usecase является потребителем этого интерфейса.
type PlantRepository interface {
	// Create сохраняет новое растение в хранилище.
	// Возвращает ID созданного растения или ошибку.
	Create(ctx context.Context, plant domain.Plant) (domain.Plant, error)

	// GetRandom извлекает указанное количество случайных растений из хранилища.
	GetRandom(ctx context.Context, count int) ([]domain.Plant, error)
}
