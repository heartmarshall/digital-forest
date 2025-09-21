package create

import (
	"context"
	"encoding/json"
	"net/http"

	domain "github.com/heartmarshall/digital-forest/backend/internal/domain/plant"
	"github.com/heartmarshall/digital-forest/backend/internal/transport/http/dto"
)

// Validator - интерфейс для валидации.
type Validator interface {
	ValidateStruct(s interface{}) map[string]string
}

// CreateUseCase - интерфейс для use case создания растения.
type CreateUseCase interface {
	Create(ctx context.Context, author, imageData string) (domain.Plant, error)
}

// CreateHandler - HTTP обработчик для создания растения.
type CreateHandler struct {
	uc        CreateUseCase
	validator Validator
}

// NewCreateHandler - конструктор для хендлера.
func NewCreateHandler(uc CreateUseCase, validator Validator) *CreateHandler {
	return &CreateHandler{
		uc:        uc,
		validator: validator,
	}
}

// CreatePlant - обработчик для POST /v1/plants
func (h *CreateHandler) CreatePlant(w http.ResponseWriter, r *http.Request) {
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
