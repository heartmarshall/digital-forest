package get_random

import (
	"context"
	"testing"
	"time"

	domain "github.com/heartmarshall/digital-forest/backend/internal/domain/plant"
	"github.com/heartmarshall/digital-forest/backend/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetRandomUseCase_GetRandom(t *testing.T) {
	tests := []struct {
		name          string
		count         int
		mockSetup     func(*testutil.MockPlantRepository)
		expectedError bool
		expectedCount int
	}{
		{
			name:  "successful get random plants",
			count: 5,
			mockSetup: func(mockRepo *testutil.MockPlantRepository) {
				expectedPlants := []domain.Plant{
					{ID: 1, Author: "author1", ImageData: "data1", CreatedAt: time.Now()},
					{ID: 2, Author: "author2", ImageData: "data2", CreatedAt: time.Now()},
					{ID: 3, Author: "author3", ImageData: "data3", CreatedAt: time.Now()},
				}
				mockRepo.On("GetRandom", mock.Anything, 5).Return(expectedPlants, nil)
			},
			expectedError: false,
			expectedCount: 3,
		},
		{
			name:  "repository error",
			count: 3,
			mockSetup: func(mockRepo *testutil.MockPlantRepository) {
				mockRepo.On("GetRandom", mock.Anything, 3).Return([]domain.Plant{}, assert.AnError)
			},
			expectedError: true,
		},
		{
			name:  "empty result",
			count: 10,
			mockSetup: func(mockRepo *testutil.MockPlantRepository) {
				mockRepo.On("GetRandom", mock.Anything, 10).Return([]domain.Plant{}, nil)
			},
			expectedError: false,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := testutil.NewMockPlantRepository()
			tt.mockSetup(mockRepo)
			useCase := NewGetRandomUseCase(mockRepo)

			// Act
			result, err := useCase.GetRandom(context.Background(), tt.count)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectedCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestNewGetRandomUseCase(t *testing.T) {
	mockRepo := testutil.NewMockPlantRepository()
	useCase := NewGetRandomUseCase(mockRepo)

	assert.NotNil(t, useCase)
	assert.Equal(t, mockRepo, useCase.repo)
}
