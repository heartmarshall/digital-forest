package usecase

import (
	"context"
	"testing"
	"time"

	domain "github.com/heartmarshall/digital-forest/backend/internal/domain/plant"
	"github.com/heartmarshall/digital-forest/backend/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPlantUseCase_Create(t *testing.T) {
	tests := []struct {
		name          string
		author        string
		imageData     string
		mockSetup     func(*testutil.MockPlantRepository)
		expectedError bool
		expectedPlant domain.Plant
	}{
		{
			name:      "successful plant creation",
			author:    "test_author",
			imageData: "base64_image_data",
			mockSetup: func(mockRepo *testutil.MockPlantRepository) {
				expectedPlant := domain.Plant{
					ID:        1,
					Author:    "test_author",
					ImageData: "base64_image_data",
					CreatedAt: time.Now().UTC(),
				}
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(plant domain.Plant) bool {
					return plant.Author == "test_author" && plant.ImageData == "base64_image_data"
				})).Return(expectedPlant, nil)
			},
			expectedError: false,
			expectedPlant: domain.Plant{
				ID:        1,
				Author:    "test_author",
				ImageData: "base64_image_data",
			},
		},
		{
			name:      "repository error",
			author:    "test_author",
			imageData: "base64_image_data",
			mockSetup: func(mockRepo *testutil.MockPlantRepository) {
				mockRepo.On("Create", mock.Anything, mock.Anything).Return(domain.Plant{}, assert.AnError)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := testutil.NewMockPlantRepository()
			tt.mockSetup(mockRepo)
			useCase := NewPlantUseCase(mockRepo)

			// Act
			result, err := useCase.Create(context.Background(), tt.author, tt.imageData)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedPlant.Author, result.Author)
				assert.Equal(t, tt.expectedPlant.ImageData, result.ImageData)
				assert.NotZero(t, result.ID)
				assert.NotZero(t, result.CreatedAt)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPlantUseCase_GetRandom(t *testing.T) {
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
			useCase := NewPlantUseCase(mockRepo)

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

func TestNewPlantUseCase(t *testing.T) {
	mockRepo := testutil.NewMockPlantRepository()
	useCase := NewPlantUseCase(mockRepo)

	assert.NotNil(t, useCase)
	assert.Equal(t, mockRepo, useCase.repo)
}
