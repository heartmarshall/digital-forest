package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	domain "github.com/heartmarshall/digital-forest/backend/internal/domain/plant"
	"github.com/heartmarshall/digital-forest/backend/internal/transport/http/dto"
)

const defaultRandomCount = 15
const maxRandomCount = 50

// Validator - интерфейс для нашей утилиты валидации.
type Validator interface {
	ValidateStruct(s interface{}) map[string]string
}

// plantUseCase - это интерфейс, который ОПРЕДЕЛЯЕТСЯ транспортным слоем (потребителем).
type plantUseCase interface {
	Create(ctx context.Context, author, imageData string) (domain.Plant, error)
	GetRandom(ctx context.Context, count int) ([]domain.Plant, error)
}

// PlantHandler - это наш HTTP обработчик.
type PlantHandler struct {
	uc        plantUseCase
	validator Validator
}

// NewPlantHandler - конструктор для хендлера.
func NewPlantHandler(uc plantUseCase, validator Validator) *PlantHandler {
	return &PlantHandler{
		uc:        uc,
		validator: validator,
	}
}

// CreatePlant - обработчик для POST /v1/plants
func (h *PlantHandler) CreatePlant(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePlantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON format"})
		return
	}

	// Выполняем автоматическую валидацию
	if validationErrors := h.validator.ValidateStruct(req); validationErrors != nil {
		respondJSON(w, http.StatusBadRequest, validationErrors)
		return
	}

	// Создаем растение через use case
	plant, err := h.uc.Create(r.Context(), req.Author, req.ImageData)
	if err != nil {
		// В реальном приложении здесь можно добавить более детальную обработку ошибок
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create plant"})
		return
	}

	// Преобразуем доменную модель в DTO для ответа
	response := dto.ToPlantResponse(plant)
	respondJSON(w, http.StatusCreated, response)
}

// GetRandomPlants - обработчик для GET /v1/plants/random
func (h *PlantHandler) GetRandomPlants(w http.ResponseWriter, r *http.Request) {
	countStr := r.URL.Query().Get("count")
	count := defaultRandomCount

	if countStr != "" {
		var err error
		count, err = strconv.Atoi(countStr)
		if err != nil || count <= 0 {
			respondJSON(w, http.StatusBadRequest, map[string]string{
				"error": "Invalid count parameter. Must be a positive integer",
			})
			return
		}

		if count > maxRandomCount {
			count = maxRandomCount
		}
	}

	plants, err := h.uc.GetRandom(r.Context(), count)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to get random plants"})
		return
	}

	responses := make([]dto.PlantResponse, len(plants))
	for i, plant := range plants {
		responses[i] = dto.ToPlantResponse(plant)
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"plants": responses,
		"count":  len(responses),
	})
}

// respondJSON - хелпер для отправки JSON-ответов.
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}
