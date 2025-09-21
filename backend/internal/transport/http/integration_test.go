package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/heartmarshall/digital-forest/backend/internal/repository/plant"
	"github.com/heartmarshall/digital-forest/backend/internal/testutil"
	"github.com/heartmarshall/digital-forest/backend/internal/transport/http/dto"
	"github.com/heartmarshall/digital-forest/backend/internal/transport/http/handlers"
	"github.com/heartmarshall/digital-forest/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockValidator для интеграционных тестов
type mockValidator struct{}

func (m *mockValidator) ValidateStruct(s interface{}) map[string]string {
	// Простая валидация для тестов
	req, ok := s.(dto.CreatePlantRequest)
	if !ok {
		return nil
	}

	errors := make(map[string]string)
	if req.Author == "" {
		errors["Author"] = "required"
	}
	if req.ImageData == "" {
		errors["ImageData"] = "required"
	}
	if len(req.Author) > 255 {
		errors["Author"] = "max length exceeded"
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

func TestHTTPIntegration(t *testing.T) {
	// Arrange
	dbPool, _, container := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, dbPool, container)

	// Setup dependencies
	plantRepo := plant.NewPlantRepo(dbPool)
	plantUseCase := usecase.NewPlantUseCase(plantRepo)
	validator := &mockValidator{}

	// Create a new chi router for testing
	router := chi.NewRouter()
	plantHandler := handlers.NewPlantHandler(plantUseCase, validator)
	router.Route("/v1/plants", func(r chi.Router) {
		r.Post("/", plantHandler.CreatePlant)
		r.Get("/random", plantHandler.GetRandomPlants)
	})

	t.Run("full API workflow", func(t *testing.T) {
		// Step 1: Create a plant
		createReq := dto.CreatePlantRequest{
			Author:    "integration_test_author",
			ImageData: "integration_test_data",
		}

		reqBody, _ := json.Marshal(createReq)
		req := httptest.NewRequest(http.MethodPost, "/v1/plants", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var createResponse dto.PlantResponse
		err := json.Unmarshal(w.Body.Bytes(), &createResponse)
		require.NoError(t, err)
		assert.NotZero(t, createResponse.ID)
		assert.Equal(t, createReq.Author, createResponse.Author)
		assert.Equal(t, createReq.ImageData, createResponse.ImageData)

		// Step 2: Create more plants for random selection
		for i := 0; i < 3; i++ {
			createReq := dto.CreatePlantRequest{
				Author:    fmt.Sprintf("author_%d", i),
				ImageData: fmt.Sprintf("data_%d", i),
			}

			reqBody, _ := json.Marshal(createReq)
			req := httptest.NewRequest(http.MethodPost, "/v1/plants", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusCreated, w.Code)
		}

		// Step 3: Get random plants
		req = httptest.NewRequest(http.MethodGet, "/v1/plants/random?count=2", nil)
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var randomResponse map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &randomResponse)
		require.NoError(t, err)

		assert.Contains(t, randomResponse, "plants")
		assert.Contains(t, randomResponse, "count")
		assert.Equal(t, float64(2), randomResponse["count"])

		plants, ok := randomResponse["plants"].([]interface{})
		require.True(t, ok)
		assert.Len(t, plants, 2)
	})

	t.Run("API error handling", func(t *testing.T) {
		// Test invalid JSON
		req := httptest.NewRequest(http.MethodPost, "/v1/plants", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Test validation error
		invalidReq := dto.CreatePlantRequest{
			Author:    "", // Empty author should fail validation
			ImageData: "valid_data",
		}

		reqBody, _ := json.Marshal(invalidReq)
		req = httptest.NewRequest(http.MethodPost, "/v1/plants", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Test invalid count parameter
		req = httptest.NewRequest(http.MethodGet, "/v1/plants/random?count=invalid", nil)
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("API with different count parameters", func(t *testing.T) {
		// Test default count
		req := httptest.NewRequest(http.MethodGet, "/v1/plants/random", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		// Test custom count
		req = httptest.NewRequest(http.MethodGet, "/v1/plants/random?count=1", nil)
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, float64(1), response["count"])

		// Test count exceeding maximum
		req = httptest.NewRequest(http.MethodGet, "/v1/plants/random?count=100", nil)
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.LessOrEqual(t, response["count"], float64(50)) // Should be capped at 50
	})
}

func TestHTTPIntegrationConcurrency(t *testing.T) {
	// Arrange
	dbPool, _, container := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, dbPool, container)

	plantRepo := plant.NewPlantRepo(dbPool)
	plantUseCase := usecase.NewPlantUseCase(plantRepo)
	validator := &mockValidator{}

	// Create a new chi router for testing
	router := chi.NewRouter()
	plantHandler := handlers.NewPlantHandler(plantUseCase, validator)
	router.Route("/v1/plants", func(r chi.Router) {
		r.Post("/", plantHandler.CreatePlant)
		r.Get("/random", plantHandler.GetRandomPlants)
	})

	t.Run("concurrent requests", func(t *testing.T) {
		// Create multiple concurrent requests
		done := make(chan bool, 10)
		errors := make(chan error, 10)

		for i := 0; i < 10; i++ {
			go func(i int) {
				createReq := dto.CreatePlantRequest{
					Author:    fmt.Sprintf("concurrent_author_%d", i),
					ImageData: fmt.Sprintf("concurrent_data_%d", i),
				}

				reqBody, _ := json.Marshal(createReq)
				req := httptest.NewRequest(http.MethodPost, "/v1/plants", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				if w.Code != http.StatusCreated {
					errors <- fmt.Errorf("unexpected status code: %d", w.Code)
					return
				}

				done <- true
			}(i)
		}

		// Wait for all requests to complete
		for i := 0; i < 10; i++ {
			select {
			case <-done:
				// Success
			case err := <-errors:
				t.Fatalf("Concurrent request failed: %v", err)
			case <-time.After(10 * time.Second):
				t.Fatal("Concurrent requests timed out")
			}
		}

		// Verify all plants were created
		req := httptest.NewRequest(http.MethodGet, "/v1/plants/random?count=15", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, response["count"], float64(10))
	})
}
