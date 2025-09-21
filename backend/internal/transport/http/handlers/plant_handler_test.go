package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	domain "github.com/heartmarshall/digital-forest/backend/internal/domain/plant"
	"github.com/heartmarshall/digital-forest/backend/internal/testutil"
	"github.com/heartmarshall/digital-forest/backend/internal/transport/http/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPlantUseCase - мок для plantUseCase
type MockPlantUseCase struct {
	mock.Mock
}

func (m *MockPlantUseCase) Create(ctx context.Context, author, imageData string) (domain.Plant, error) {
	args := m.Called(ctx, author, imageData)
	return args.Get(0).(domain.Plant), args.Error(1)
}

func (m *MockPlantUseCase) GetRandom(ctx context.Context, count int) ([]domain.Plant, error) {
	args := m.Called(ctx, count)
	return args.Get(0).([]domain.Plant), args.Error(1)
}

func TestPlantHandler_CreatePlant(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockPlantUseCase, *testutil.MockValidator)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "successful plant creation",
			requestBody: dto.CreatePlantRequest{
				Author:    "test_author",
				ImageData: "base64_image_data",
			},
			mockSetup: func(mockUC *MockPlantUseCase, mockValidator *testutil.MockValidator) {
				mockValidator.On("ValidateStruct", mock.Anything).Return(nil)
				expectedPlant := domain.Plant{
					ID:        1,
					Author:    "test_author",
					ImageData: "base64_image_data",
					CreatedAt: time.Now().UTC(),
				}
				mockUC.On("Create", mock.Anything, "test_author", "base64_image_data").Return(expectedPlant, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name:        "invalid JSON",
			requestBody: "invalid json",
			mockSetup: func(mockUC *MockPlantUseCase, mockValidator *testutil.MockValidator) {
				// No mocks needed for invalid JSON
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "validation error",
			requestBody: dto.CreatePlantRequest{
				Author:    "", // Empty author should fail validation
				ImageData: "base64_image_data",
			},
			mockSetup: func(mockUC *MockPlantUseCase, mockValidator *testutil.MockValidator) {
				validationErrors := map[string]string{
					"Author": "required",
				}
				mockValidator.On("ValidateStruct", mock.Anything).Return(validationErrors)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "use case error",
			requestBody: dto.CreatePlantRequest{
				Author:    "test_author",
				ImageData: "base64_image_data",
			},
			mockSetup: func(mockUC *MockPlantUseCase, mockValidator *testutil.MockValidator) {
				mockValidator.On("ValidateStruct", mock.Anything).Return(nil)
				mockUC.On("Create", mock.Anything, "test_author", "base64_image_data").Return(domain.Plant{}, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockUC := &MockPlantUseCase{}
			mockValidator := testutil.NewMockValidator()
			tt.mockSetup(mockUC, mockValidator)

			handler := NewPlantHandler(mockUC, mockValidator)

			var reqBody []byte
			if tt.requestBody != nil {
				reqBody, _ = json.Marshal(tt.requestBody)
			} else {
				reqBody = []byte("invalid json")
			}

			req := httptest.NewRequest(http.MethodPost, "/v1/plants", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Act
			handler.CreatePlant(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				// For validation errors, the response contains field-specific errors
				// For other errors, it contains an "error" field
				if tt.name == "validation error" {
					assert.Contains(t, response, "Author")
				} else {
					assert.Contains(t, response, "error")
				}
			} else {
				var response dto.PlantResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotZero(t, response.ID)
				assert.Equal(t, "test_author", response.Author)
				assert.Equal(t, "base64_image_data", response.ImageData)
			}

			mockUC.AssertExpectations(t)
			mockValidator.AssertExpectations(t)
		})
	}
}

func TestPlantHandler_GetRandomPlants(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		mockSetup      func(*MockPlantUseCase)
		expectedStatus int
		expectedCount  int
		expectedError  bool
	}{
		{
			name:        "successful get with default count",
			queryParams: "",
			mockSetup: func(mockUC *MockPlantUseCase) {
				expectedPlants := []domain.Plant{
					{ID: 1, Author: "author1", ImageData: "data1", CreatedAt: time.Now()},
					{ID: 2, Author: "author2", ImageData: "data2", CreatedAt: time.Now()},
				}
				mockUC.On("GetRandom", mock.Anything, 15).Return(expectedPlants, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
			expectedError:  false,
		},
		{
			name:        "successful get with custom count",
			queryParams: "?count=5",
			mockSetup: func(mockUC *MockPlantUseCase) {
				expectedPlants := []domain.Plant{
					{ID: 1, Author: "author1", ImageData: "data1", CreatedAt: time.Now()},
				}
				mockUC.On("GetRandom", mock.Anything, 5).Return(expectedPlants, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
			expectedError:  false,
		},
		{
			name:        "invalid count parameter",
			queryParams: "?count=invalid",
			mockSetup: func(mockUC *MockPlantUseCase) {
				// No mocks needed for invalid count
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:        "negative count parameter",
			queryParams: "?count=-1",
			mockSetup: func(mockUC *MockPlantUseCase) {
				// No mocks needed for negative count
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:        "count exceeds maximum",
			queryParams: "?count=100",
			mockSetup: func(mockUC *MockPlantUseCase) {
				expectedPlants := []domain.Plant{
					{ID: 1, Author: "author1", ImageData: "data1", CreatedAt: time.Now()},
				}
				mockUC.On("GetRandom", mock.Anything, 50).Return(expectedPlants, nil) // Should be capped at 50
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
			expectedError:  false,
		},
		{
			name:        "use case error",
			queryParams: "?count=5",
			mockSetup: func(mockUC *MockPlantUseCase) {
				mockUC.On("GetRandom", mock.Anything, 5).Return([]domain.Plant{}, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockUC := &MockPlantUseCase{}
			mockValidator := testutil.NewMockValidator()
			tt.mockSetup(mockUC)

			handler := NewPlantHandler(mockUC, mockValidator)

			req := httptest.NewRequest(http.MethodGet, "/v1/plants/random"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			// Act
			handler.GetRandomPlants(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			} else {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "plants")
				assert.Contains(t, response, "count")
				assert.Equal(t, float64(tt.expectedCount), response["count"])
			}

			mockUC.AssertExpectations(t)
		})
	}
}

func TestNewPlantHandler(t *testing.T) {
	mockUC := &MockPlantUseCase{}
	mockValidator := testutil.NewMockValidator()

	handler := NewPlantHandler(mockUC, mockValidator)

	assert.NotNil(t, handler)
	assert.Equal(t, mockUC, handler.uc)
	assert.Equal(t, mockValidator, handler.validator)
}
