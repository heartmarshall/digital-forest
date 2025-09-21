package create

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

// MockCreateUseCase - мок для CreateUseCase
type MockCreateUseCase struct {
	mock.Mock
}

func (m *MockCreateUseCase) Create(ctx context.Context, author, imageData string) (domain.Plant, error) {
	args := m.Called(ctx, author, imageData)
	return args.Get(0).(domain.Plant), args.Error(1)
}

func TestCreateHandler_CreatePlant(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockCreateUseCase, *testutil.MockValidator)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "successful plant creation",
			requestBody: dto.CreatePlantRequest{
				Author:    "test_author",
				ImageData: "base64_image_data",
			},
			mockSetup: func(mockUC *MockCreateUseCase, mockValidator *testutil.MockValidator) {
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
			mockSetup: func(mockUC *MockCreateUseCase, mockValidator *testutil.MockValidator) {
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
			mockSetup: func(mockUC *MockCreateUseCase, mockValidator *testutil.MockValidator) {
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
			mockSetup: func(mockUC *MockCreateUseCase, mockValidator *testutil.MockValidator) {
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
			mockUC := &MockCreateUseCase{}
			mockValidator := testutil.NewMockValidator()
			tt.mockSetup(mockUC, mockValidator)

			handler := NewCreateHandler(mockUC, mockValidator)

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

func TestNewCreateHandler(t *testing.T) {
	mockUC := &MockCreateUseCase{}
	mockValidator := testutil.NewMockValidator()

	handler := NewCreateHandler(mockUC, mockValidator)

	assert.NotNil(t, handler)
	assert.Equal(t, mockUC, handler.uc)
	assert.Equal(t, mockValidator, handler.validator)
}
