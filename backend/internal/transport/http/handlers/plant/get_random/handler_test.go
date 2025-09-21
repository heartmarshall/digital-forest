package get_random

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	domain "github.com/heartmarshall/digital-forest/backend/internal/domain/plant"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGetRandomUseCase - мок для GetRandomUseCase
type MockGetRandomUseCase struct {
	mock.Mock
}

func (m *MockGetRandomUseCase) GetRandom(ctx context.Context, count int) ([]domain.Plant, error) {
	args := m.Called(ctx, count)
	return args.Get(0).([]domain.Plant), args.Error(1)
}

func TestGetRandomHandler_GetRandomPlants(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		mockSetup      func(*MockGetRandomUseCase)
		expectedStatus int
		expectedCount  int
		expectedError  bool
	}{
		{
			name:        "successful get with default count",
			queryParams: "",
			mockSetup: func(mockUC *MockGetRandomUseCase) {
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
			mockSetup: func(mockUC *MockGetRandomUseCase) {
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
			mockSetup: func(mockUC *MockGetRandomUseCase) {
				// No mocks needed for invalid count
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:        "negative count parameter",
			queryParams: "?count=-1",
			mockSetup: func(mockUC *MockGetRandomUseCase) {
				// No mocks needed for negative count
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:        "count exceeds maximum",
			queryParams: "?count=100",
			mockSetup: func(mockUC *MockGetRandomUseCase) {
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
			mockSetup: func(mockUC *MockGetRandomUseCase) {
				mockUC.On("GetRandom", mock.Anything, 5).Return([]domain.Plant{}, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockUC := &MockGetRandomUseCase{}
			tt.mockSetup(mockUC)

			handler := NewGetRandomHandler(mockUC)

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

func TestNewGetRandomHandler(t *testing.T) {
	mockUC := &MockGetRandomUseCase{}

	handler := NewGetRandomHandler(mockUC)

	assert.NotNil(t, handler)
	assert.Equal(t, mockUC, handler.uc)
}
