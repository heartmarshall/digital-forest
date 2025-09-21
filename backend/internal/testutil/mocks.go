package testutil

import (
	"context"

	domain "github.com/heartmarshall/digital-forest/backend/internal/domain/plant"
	"github.com/stretchr/testify/mock"
)

// MockPlantRepository - мок для PlantRepository
type MockPlantRepository struct {
	mock.Mock
}

func (m *MockPlantRepository) Create(ctx context.Context, plant domain.Plant) (domain.Plant, error) {
	args := m.Called(ctx, plant)
	return args.Get(0).(domain.Plant), args.Error(1)
}

func (m *MockPlantRepository) GetRandom(ctx context.Context, count int) ([]domain.Plant, error) {
	args := m.Called(ctx, count)
	return args.Get(0).([]domain.Plant), args.Error(1)
}

// MockValidator - мок для валидатора
type MockValidator struct {
	mock.Mock
}

func (m *MockValidator) ValidateStruct(s interface{}) map[string]string {
	args := m.Called(s)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(map[string]string)
}

// NewMockPlantRepository создает новый мок репозитория
func NewMockPlantRepository() *MockPlantRepository {
	return &MockPlantRepository{}
}

// NewMockValidator создает новый мок валидатора
func NewMockValidator() *MockValidator {
	return &MockValidator{}
}

// PlantRepositoryInterface определяет интерфейс для репозитория растений
type PlantRepositoryInterface interface {
	Create(ctx context.Context, plant domain.Plant) (domain.Plant, error)
	GetRandom(ctx context.Context, count int) ([]domain.Plant, error)
}

// AssertPlantRepositoryInterface проверяет, что мок реализует интерфейс
var _ PlantRepositoryInterface = (*MockPlantRepository)(nil)
