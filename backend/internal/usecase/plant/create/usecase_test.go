package create

import (
	"context"
	"testing"
	"time"

	domain "github.com/heartmarshall/digital-forest/backend/internal/domain/plant"
	"github.com/heartmarshall/digital-forest/backend/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUseCase_Create(t *testing.T) {
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
			useCase := NewCreateUseCase(mockRepo)

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

func TestNewCreateUseCase(t *testing.T) {
	mockRepo := testutil.NewMockPlantRepository()
	useCase := NewCreateUseCase(mockRepo)

	assert.NotNil(t, useCase)
	assert.Equal(t, mockRepo, useCase.repo)
}
